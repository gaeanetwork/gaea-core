package task

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gaeanetwork/gaea-core/services/transmission"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/system"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/container/util"
	"github.com/pkg/errors"
)

// ExecuteRequest is to execute a task through this request
type ExecuteRequest struct {
	TaskID        string `form:"task_id"`
	AlgorithmID   string `form:"algorithm_id"`
	Executor      string `form:"executor"`
	ContainerType string `form:"container_type"`
	Hash          string `form:"hash"`
	Signature     string `form:"signature"`
}

// Execute executes a tee task by sending a execute transaction to the blockchain.
func Execute(req *ExecuteRequest) (string, error) {
	containerName, err := GetContainerName(ChaincodeName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get container name, chaincode: %s", ChaincodeName)
	}

	teetask, err := GetByID(req.TaskID)
	if err != nil {
		return "", errors.Wrapf(err, "Tee task cannot found, taskID: %v", req.TaskID)
	}

	if err = uploadToTeetaskContainer(containerName, req.AlgorithmID, teetask); err != nil {
		return "", errors.Wrapf(err, "failed to upload to tee task container, address: %v", teetask.ResultAddress)
	}

	service, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get smart contract service, platform: %s", tee.ImplementPlatform)
	}

	executeID, err := service.Invoke(ChaincodeName, []string{MethodExecute, req.TaskID, req.Executor, req.ContainerType, req.Hash, req.Signature})
	if err != nil {
		return "", fmt.Errorf("Error executing tee task, error: %v", err)
	}

	if err = downloadFromTeetaskContainer(containerName, executeID, teetask); err != nil {
		return "", errors.Wrapf(err, "failed to downlaod from tee task container, address: %v", teetask.ResultAddress)
	}

	return string(executeID), nil
}

// GetContainerName get chaincode docker container name by chaincode Name
func GetContainerName(chaincodeName string) (string, error) {
	config, err := chaincode.GetConfig(chaincodeName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get chaincode config, name: %s", chaincodeName)
	}

	definition, err := system.GetChaincodeDefinition(config.ChannelID, config.ChaincodeName)
	if err != nil {
		return "", errors.Wrap(err, "failed to get chaincode definition")
	}

	return chaincode.GetContainerName(definition.CCName(), definition.CCVersion()), nil
}

func uploadToTeetaskContainer(containerName, algorithmID string, teetask *tee.Task) error {
	// Get data and algorithm byte slice
	filesToUpload, err := getDataToUpload(teetask)
	if err != nil {
		return errors.Wrap(err, "failed to get files to upload")
	}

	// TODO - userDir := filepath.Join(transmission.DefaultLocation, algorithmData.Owner, algorithmID)
	fileDir := filepath.Join(transmission.DefaultLocation, algorithmID)
	algorithm, err := ioutil.ReadFile(fileDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read algorithm file, path: %s", fileDir)
	}
	filesToUpload[filepath.Join(teetask.ResultAddress, AlgorithmName)] = algorithm

	return errors.Wrapf(upload(containerName, filesToUpload), "containerName: %s, teetask.ResultAddress: %s", containerName, teetask.ResultAddress)
}

func getDataToUpload(teetask *tee.Task) (map[string][]byte, error) {
	filesToUpload := make(map[string][]byte)
	for _, notificationID := range teetask.DataNotifications {
		notification, err := tee.GetNotification(notificationID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get notification, id: %s", notificationID)
		}

		if notification.Status != tee.Authorized {
			return nil, fmt.Errorf("Failed to check request data status, notificationID: %s, status: %s, message: %s", notificationID, notification.Status, notification.RefusedReason)
		}

		if notification.DataInfo.DataStoreType == tee.Local {
			fileDir := filepath.Join(transmission.DefaultLocation, notification.DataInfo.DataStoreAddress)
			data, err := ioutil.ReadFile(fileDir)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to read file, path: %s", fileDir)
			}

			filesToUpload[filepath.Join(teetask.ResultAddress, notification.DataInfo.DataStoreAddress)] = data
		}
	}

	return filesToUpload, nil
}

func upload(containerName string, filesToUpload map[string][]byte) error {
	if len(filesToUpload) == 0 {
		return nil
	}

	// the docker upload API takes a tar file, so we need to first
	// consolidate the file entries to a tar
	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	for path, fileToUpload := range filesToUpload {
		util.WriteBytesToPackage(path, fileToUpload, tw)
	}

	// Write the tar file out
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Error writing files to upload to Docker instance into a temporary tar blob: %s", err)
	}
	gw.Close()

	client, err := util.NewDockerClient()
	if err != nil {
		return fmt.Errorf("Error getting docker client, error: %v", err)
	}

	return client.UploadToContainer(containerName, docker.UploadToContainerOptions{
		InputStream:          bytes.NewReader(payload.Bytes()),
		Path:                 "/",
		NoOverwriteDirNonDir: false,
	})
}

func downloadFromTeetaskContainer(containerName string, executeID []byte, teetask *tee.Task) error {
	for _, notificationID := range teetask.DataNotifications {
		notification, err := tee.GetNotification(notificationID)
		if err != nil {
			return errors.Wrapf(err, "failed to get notification, id: %s", notificationID)
		}

		if notification.DataInfo.DataStoreType == tee.Local {
			client, err := util.NewDockerClient()
			if err != nil {
				return fmt.Errorf("Error getting docker client, error: %v", err)
			}

			buffer, logPath := bytes.NewBuffer(nil), filepath.Join(teetask.ResultAddress, notification.Data.Owner+".log")
			err = client.DownloadFromContainer(containerName, docker.DownloadFromContainerOptions{
				OutputStream: buffer,
				Path:         logPath,
			})
			if err != nil {
				return fmt.Errorf("Error donwloading execution.log from the container instance %s: %s", containerName, err)
			}

			reader := tar.NewReader(buffer)
			if _, err = reader.Next(); err != nil {
				return errors.Wrapf(err, "failed to read output stream header")
			}

			var result bytes.Buffer
			for {
				data := make([]byte, 1024)
				n, err := reader.Read(data)
				if err == io.EOF {
					result.Write(data[:n])
					break
				}

				if err != nil {
					return errors.Wrapf(err, "failed to read execution log")
				}

				result.Write(data[:n])
			}

			fileID := sha256.Sum256(append(executeID, []byte(notification.Data.Owner)...))
			if err = ioutil.WriteFile(filepath.Join(transmission.DefaultLocation, hex.EncodeToString(fileID[:])), result.Bytes(), os.ModePerm); err != nil {
				return errors.Wrap(err, "failed to write log file")
			}
		}
	}

	return nil
}
