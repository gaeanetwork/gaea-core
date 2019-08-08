package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gaeanetwork/gaea-core/tee"
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
