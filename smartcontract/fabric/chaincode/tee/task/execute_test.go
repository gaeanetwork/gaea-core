package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/gaeanetwork/gaea-core/tee/task"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_execute(t *testing.T) {
	stub := getTeeTaskMockStub()

	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodExecute)})
	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte("")})
	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte(""), []byte("")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte(""), []byte("executor"), []byte("container")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte("taskID"), []byte(""), []byte("executor"), []byte("container")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte("taskID"), []byte(""), []byte("container")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodExecute), []byte("taskID"), []byte("executor"), []byte("")})
	// create a teetask
	taskID := createTeeTask(t, "1", stub)
	ioutil.WriteFile("/tmp/teetask/test/1/main", []byte("Testing"), os.ModePerm)
	defer os.Remove("/tmp/teetask/test/1/main")

	// Invalid taskID
	resp := stub.MockInvoke("2", [][]byte{[]byte(task.MethodExecute), []byte("taskID"), []byte("executor"), []byte("0")})
	assert.Equal(t, shim.ERROR, int(resp.Status))
	assert.Contains(t, resp.Message, "Data does not exist")

	// Invalid notification status
	resp2 := stub.MockInvoke("4", [][]byte{[]byte(task.MethodExecute), taskID, []byte(PubHexForTest), []byte("0")})
	assert.Equal(t, shim.ERROR, int(resp2.Status))
	assert.Contains(t, resp2.Message, "Failed to check request data status")
}

func Test_executeSignature(t *testing.T) {
	stub := getTeeTaskMockStub()

	taskID := createTeeTask(t, "1", stub)
	ioutil.WriteFile("/tmp/teetask/test/1/main", []byte("Testing"), os.ModePerm)
	defer os.Remove("/tmp/teetask/test/1/main")

	response := stub.MockInvoke("2", [][]byte{[]byte("get"), taskID})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	var teetask tee.Task
	err := json.Unmarshal(response.Payload, &teetask)
	assert.NoError(t, err)
	assert.Equal(t, taskID, []byte(teetask.ID))
	_, ok := teetask.Partners[PubHexForTest]
	assert.True(t, ok)

	// Add sigs
	args := []string{string(taskID), PubHexForTest, "0"}
	privBytes, err := hex.DecodeString(PrivHexForTest)
	assert.NoError(t, err)
	hash, sigs, err := chaincode.GetArgsHashAndSignatures(privBytes, args)
	assert.NoError(t, err)

	resp := stub.MockInvoke("3", [][]byte{[]byte(task.MethodExecute), taskID, []byte(PubHexForTest), []byte("0"), []byte(hex.EncodeToString(hash)), sigs})
	assert.Equal(t, shim.ERROR, int(resp.Status))
	assert.Contains(t, resp.Message, "Failed to check request data status")

	// Invalid hash
	resp1 := stub.MockInvoke("4", [][]byte{[]byte(task.MethodExecute), taskID, []byte(PubHexForTest), []byte("0"), []byte("hex.EncodeToString(hash)"), sigs})
	assert.Equal(t, shim.ERROR, int(resp1.Status))
	assert.Contains(t, resp1.Message, "Failed to verify ecdsa signature")

	// Invalid signature
	sigs1 := []string{"signature"}
	sigsBytes1, err := json.Marshal(sigs1)
	assert.NoError(t, err)

	resp2 := stub.MockInvoke("5", [][]byte{[]byte(task.MethodExecute), taskID, []byte(PubHexForTest), []byte("0"), []byte("hex.EncodeToString(hash)"), sigsBytes1})
	assert.Equal(t, shim.ERROR, int(resp2.Status))
	assert.Contains(t, resp2.Message, "Failed to verify ecdsa signature")
}

func Test_DockerVolume(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	assert.NoError(t, err)
	imgs, err := dockerClient.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		if img.RepoTags[0] == "dev1-peer1-teetask-1.7-33290189c29e65cc17c60a0db90fdb3fcb964b0803426953a9bd993e64bb545d:latest" {
			volume, err := dockerClient.CreateVolume(docker.CreateVolumeOptions{
				Name:   "tardis",
				Driver: "local",
			})
			assert.NoError(t, err)

			volume1, err1 := dockerClient.InspectVolume(volume.Name)
			assert.NoError(t, err1)
			assert.Equal(t, volume, volume1)
		}
	}
}

func Test_checkNotificationAuthorized(t *testing.T) {
	dataInfo := &tee.DataInfo{}
	notification := &tee.Notification{Status: tee.Authorized, DataInfo: dataInfo}
	result, err := json.Marshal(notification)
	assert.NoError(t, err)
	stub := getTeeTaskMockStubByTeeController(result, "")

	notification1, err := checkNotificationAuthorized(stub, "1")
	assert.NoError(t, err)
	assert.Equal(t, dataInfo, notification1.DataInfo)

	// Not Authorized
	notification = &tee.Notification{Status: tee.Refused, DataInfo: dataInfo}
	result, err = json.Marshal(notification)
	assert.NoError(t, err)
	stub = getTeeTaskMockStubByTeeController(result, "")
	_, err = checkNotificationAuthorized(stub, "1")
	assert.Contains(t, err.Error(), "Failed to check request data status")

	// Not notification
	stub = getTeeTaskMockStubByTeeController([]byte("result"), "")
	_, err = checkNotificationAuthorized(stub, "1")
	assert.Contains(t, err.Error(), "invalid character")

	// Invoke tee chaincode not ok
	stub = shim.NewMockStub(task.ChaincodeName, new(TeeTaskChaincode))
	assert.Panics(t, func() { checkNotificationAuthorized(stub, "1") })
}
