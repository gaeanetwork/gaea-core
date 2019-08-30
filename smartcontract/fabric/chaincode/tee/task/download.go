package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/pipeline"
	"gitlab.com/jaderabbit/go-rabbit/tee/storage/azurehelper"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

// DownloadData if the the request status is authorized, download the data
func downloadData(stub shim.ChaincodeStubInterface, teetask *tee.Task) ([][]byte, error) {
	type Result struct {
		Err  error
		Data []byte
	}

	download := func(done chan interface{}, notificationStream <-chan interface{}) <-chan interface{} {
		dataStream := make(chan interface{})
		go func() {
			defer close(dataStream)
			for notificationID := range notificationStream {
				select {
				case <-done:
					return
				default:
				}
				// Check that the notification status is authorized, and if so, return the notification's dataInfo.
				go func() { logger.Infof("start to check Notification by notificationID: %v", notificationID) }()
				notification, err := checkNotificationAuthorized(stub, notificationID.(string))
				if err != nil {
					dataStream <- Result{Err: errors.Wrapf(err, "failed to check notification authorized status, notificationsID: %s", notificationID)}
					return
				}
				go func() { logger.Infof("checked Notification by notificationID: %v", notificationID) }()

				// Decrypt data address
				dataAddress, err := decryptDataAddress(stub, notification.DataInfo)
				if err != nil {
					dataStream <- Result{Err: errors.Wrapf(err, "failed to decrypt ciphertext data address, notification.DataInfo: %v", notification.DataInfo)}
					return
				}

				// Download data from data address
				aesKey, err := decryptEncryptedKey(stub, notification.DataInfo.EncryptedKey)
				if err != nil {
					dataStream <- Result{Err: errors.Wrapf(err, "failed to generate shared secret key, pubKey: %s", notification.DataInfo.EncryptedKey)}
					return
				}
				// log.Println("aesKey:", hex.EncodeToString(aesKey))

				ciphertextData, err := downloadDataFromAddress(teetask.ResultAddress, dataAddress, notification, aesKey)
				if err != nil {
					dataStream <- Result{Err: errors.Wrapf(err, "failed to download data, resultAddress: %s, address: %s, type: %d", teetask.ResultAddress, dataAddress, notification.DataInfo.DataStoreType)}
					return
				}

				// Decrypt data
				plaintextData, err := decryptData(stub, ciphertextData, notification.DataInfo)
				if err != nil {
					dataStream <- Result{Err: errors.Wrapf(err, "failed to decrypt data, ciphertext data: %s", ciphertextData)}
					return
				}

				dataStream <- Result{Data: plaintextData}
			}
		}()

		return dataStream
	}

	done := make(chan interface{})
	defer close(done)

	notifications := make([]interface{}, 0)
	for _, notificationID := range teetask.DataNotifications {
		notifications = append(notifications, notificationID)
	}

	dataList := make([][]byte, 0)
	for result := range download(done, pipeline.Generator(done, notifications...)) {
		r := result.(Result)
		if r.Err != nil {
			return nil, r.Err
		}

		dataList = append(dataList, r.Data)
	}

	logger.Infof("txID: %s, dataList length: %d\n", stub.GetTxID(), len(dataList))
	if len(dataList) == 0 {
		return nil, fmt.Errorf("Error downloading the data list is empty")
	}

	return dataList, nil
}

// download data from data address
func downloadDataFromAddress(resultAddress, dataAddress string, notification *tee.Notification, aesKey []byte) (data []byte, err error) {
	switch notification.DataInfo.DataStoreType {
	case tee.Local:
		hooks[notification.Data.Owner] = func(data []byte) error {
			log.Println("local uploading...")
			ciphertext, err := crypto.AesEncrypt(data, aesKey)
			if err != nil {
				return errors.Wrapf(err, "failed to encrypt data, owner: %s", notification.Data.Owner)
			}

			logPath := filepath.Join(resultAddress, notification.Data.Owner+".log")
			log.Println("write execution.log in:", logPath)
			return ioutil.WriteFile(logPath, []byte(hex.EncodeToString(ciphertext)), 0755)
		}

		return ioutil.ReadFile(filepath.Join(resultAddress, dataAddress))
	case tee.Azure:
		// TODO - handle the account and access key
		os.Setenv(azurehelper.AzureStorageAccountEnvKey, "azureteeaccount")
		os.Setenv(azurehelper.AzureStorageAccessKeyEnvKey, "0jlvPT+Gw2+Y4ltGfuOXCkw91QQI82gsL2RjHbQPCmo7VlDneujTFnu+B7a/FC5tfVdCkVCZRti1Dpfw0Evaaw==")

		dataAddresses := strings.Split(dataAddress, task.AzureSplitSep)
		if len(dataAddresses) != 2 {
			return nil, fmt.Errorf("Invalid data store address, azure data store address should be `containerName~filePath`, now: %s", dataAddress)
		}
		containerName, filePath := dataAddresses[0], dataAddresses[1]

		hooks[notification.Data.Owner] = func(data []byte) error {
			log.Println("azure uploading...")
			ciphertext, err := crypto.AesEncrypt(data, aesKey)
			if err != nil {
				return errors.Wrapf(err, "failed tp encrypt data, owner: %s", notification.Data.Owner)
			}

			logPath := filepath.Join(resultAddress, notification.Data.Owner+".log")
			log.Println("write execution.log in:", logPath)
			return azurehelper.UploadFileToContainer(containerName, logPath, []byte(hex.EncodeToString(ciphertext)))
		}

		return azurehelper.DownloadFilesFromContainer(containerName, filePath)
	default:
		return nil, fmt.Errorf("Unimplemented data store type, type: %d", notification.DataInfo.DataStoreType)
	}
}
