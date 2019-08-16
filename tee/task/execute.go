package task

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/container/util"
	"github.com/pkg/errors"
)

// ExecuteRequest is to execute a task through this request
type ExecuteRequest struct {
	TaskID        string `form:"task_id"`
	Algorithm     []byte `form:"algorithm"`
	Executor      string `form:"executor"`
	ContainerType string `form:"container_type"`
	Hash          string `form:"hash"`
	Signature     string `form:"signature"`
}

func uploadToTeetaskContainer(containerName string, algorithm []byte, teetask *tee.Task) error {
	// Get data and algorithm byte slice
	filesToUpload, err := getDataToUpload(teetask)
	if err != nil {
		return errors.Wrap(err, "failed to get files to upload")
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
			err := filepath.Walk(teetask.ResultAddress, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}

				data, err := ioutil.ReadFile(path)
				if err != nil {
					return errors.Wrapf(err, "failed to read file, path: %s", path)
				}

				filesToUpload[path] = data
				return nil
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to walk folder, folder: %s", teetask.ResultAddress)
			}

			break
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

func downloadFromTeetaskContainer(containerName string, teetask *tee.Task) error {
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

			if err = ioutil.WriteFile(logPath, result.Bytes(), os.ModePerm); err != nil {
				return errors.Wrap(err, "failed to write log file")
			}
		}
	}

	return nil
}
