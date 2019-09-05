package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/gaeanetwork/gaea-core/tee/task"
	"github.com/stretchr/testify/assert"
)

var (
	containerNameForTests = "tee-container"
	filePathForTests      = "teetest/data/A.secret"
	resultAddress         = os.TempDir()
	dataAddress           = "data.go"
)

func Test_downloadData(t *testing.T) {
	filePath := filepath.Join(resultAddress, dataAddress)
	err := ioutil.WriteFile(filePath, plaintextForTests, os.ModePerm)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	dataInfo := &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedType: tee.UnEncrypted}
	notification := &tee.Notification{ID: "1", Data: &tee.SharedData{Owner: PubHexForTest}, Status: tee.Authorized, DataInfo: dataInfo}
	result, err := json.Marshal(notification)
	assert.NoError(t, err)
	stub := getTeeTaskMockStubByTeeController(result, "")
	task := &tee.Task{DataNotifications: map[string]string{"?": notification.ID}, ResultAddress: resultAddress}

	dataList, err := downloadData(stub, task)
	assert.NoError(t, err)
	assert.Len(t, dataList, 1)
	assert.Equal(t, plaintextForTests, dataList[0])
}

func Test_downloadData_Error(t *testing.T) {
	// Failed to decrypt data address
	dataInfo := &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedType: tee.EncryptedType(23)}
	notification := &tee.Notification{ID: "1", Status: tee.Authorized, DataInfo: dataInfo}
	result, err := json.Marshal(notification)
	assert.NoError(t, err)
	stub := getTeeTaskMockStubByTeeController(result, "")
	task := &tee.Task{DataNotifications: map[string]string{"?": notification.ID}, ResultAddress: resultAddress}

	_, err = downloadData(stub, task)
	assert.Contains(t, err.Error(), "failed to decrypt ciphertext data address")

	// Failed to download data from address
	dataInfo = &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedType: tee.UnEncrypted}
	notification = &tee.Notification{ID: "1", Data: &tee.SharedData{Owner: PubHexForTest}, Status: tee.Authorized, DataInfo: dataInfo}
	result, err = json.Marshal(notification)
	assert.NoError(t, err)
	stub = getTeeTaskMockStubByTeeController(result, "")

	_, err = downloadData(stub, task)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func Test_downloadDataFromAddress(t *testing.T) {
	notification := &tee.Notification{Data: &tee.SharedData{Owner: PubHexForTest}, DataInfo: &tee.DataInfo{}}
	// Local
	notification.DataInfo.DataStoreType = tee.Local
	filePath := filepath.Join(resultAddress, dataAddress)
	err := ioutil.WriteFile(filePath, plaintextForTests, os.ModePerm)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	data, err := downloadDataFromAddress(resultAddress, dataAddress, notification, []byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, plaintextForTests, data)
	testUploadHookInLocal(t, data)

	// clean local hook
	delete(hooks, PubHexForTest)
	_, exists := hooks[PubHexForTest]
	assert.False(t, exists)

	// Azure
	notification.DataInfo.DataStoreType = tee.Azure
	dataAddress = strings.Join([]string{containerNameForTests, filePathForTests}, task.AzureSplitSep)
	data, err = downloadDataFromAddress(resultAddress, dataAddress, notification, []byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, string(data), "0xe766a23e17883fd67588264febe2c3a5ec7ebdda1e3b1612772b17d7d2c9cb496207ed9f741e8048c91b44e3b371d1527188765043e4ed51c6856b2d6045cbfe0d8892e7be7d49665aebc598c4de9495")
	testUploadHookInAzure(t, data)
}

func Test_downloadDataFromAddress_Error(t *testing.T) {
	// Invalid data store type
	_, err := downloadDataFromAddress(resultAddress, dataAddress, &tee.Notification{DataInfo: &tee.DataInfo{DataStoreType: tee.DataStoreType(823)}}, []byte("123"))
	assert.Contains(t, err.Error(), "Unimplemented data store type")

	// Invalid data split sep for azure
	dataAddress = strings.Join([]string{containerNameForTests, filePathForTests}, "task.AzureSplitSep")
	_, err = downloadDataFromAddress("", dataAddress, &tee.Notification{DataInfo: &tee.DataInfo{DataStoreType: tee.Azure}}, []byte("123"))
	assert.Contains(t, err.Error(), "azure data store address should be `containerName~filePath`")
}
