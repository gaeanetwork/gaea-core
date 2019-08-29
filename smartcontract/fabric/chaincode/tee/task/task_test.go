package main

/*
    ============================ Deploy teetask chaincode ============================
	./peer chaincode package teetaskpack.out -n tee_exec -v 1.0 -s -S -p gitlab.com/jaderabbit/go-rabbit/chaincode/tee/task
	mkdir $HOME/chaincodes/tee
	mv -fv teetaskpack.out $HOME/chaincodes/tee/teetaskpack.out

	# install
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install $HOME/chaincodes/tee/teetaskpack.out

	# instantiate
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode instantiate -C syschannel -n tee_exec -v 1.0 -c '{"Args":["0x307702010104204368376222802d1a941f2eb0b7186a2c75f75e368946f923ad37e7c7718c2d7aa00a06082a8648ce3d030107a14403420004e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01", "0x04e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"]}' -o orderer.rabbit.com:7050

	# upgrade
	CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode upgrade -C syschannel -n tee_exec -v 1.1 -c '{"Args":["0x307702010104204368376222802d1a941f2eb0b7186a2c75f75e368946f923ad37e7c7718c2d7aa00a06082a8648ce3d030107a14403420004e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01", "0x04e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"]}' -o orderer.rabbit.com:7050

	============================ Test tee chaincode ============================
	./peer chaincode invoke -C syschannel -n tee_exec -c '{"Args":["create","taskID","0x1111"]}' -o orderer.rabbit.com:7050
	./peer chaincode invoke -C syschannel -n tee_exec -c '{"Args":["getall"]}' -o orderer.rabbit.com:7050
	./peer chaincode invoke -C syschannel -n tee_exec -c '{"Args":["get","taskID"]}' -o orderer.rabbit.com:7050
*/

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

func Test_Init(t *testing.T) {
	stub := shim.NewMockStub("teetask", new(TeeTaskChaincode))

	args := [][]byte{[]byte("privHex"), []byte("pubHex")}
	response := stub.MockInit("1", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	// Invalid parameters
	args1 := [][]byte{[]byte(""), []byte("")}
	response1 := stub.MockInit("2", args1)
	assert.Equal(t, shim.ERROR, int(response1.Status))
	assert.Contains(t, response1.Message, "argument must be a non-empty string")

	// Invalid parameters length
	args2 := [][]byte{}
	response2 := stub.MockInit("2", args2)
	assert.Equal(t, shim.ERROR, int(response2.Status))
	assert.Contains(t, response2.Message, "Incorrect number of arguments.")
}

func filedEmptyError(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("1", args)
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.Contains(t, response.Message, "argument must be a non-empty string")
}

func IncorrectNumberArgsError(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("1", args)
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.Contains(t, response.Message, "Incorrect number of arguments.")
}

func Test_get(t *testing.T) {
	stub := getTeeTaskMockStub()

	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodGet)})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodGet), []byte("")})

	taskID := createTeeTask(t, "1", stub)

	response := stub.MockInvoke("2", [][]byte{[]byte(task.MethodGet), taskID})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	var teetask tee.Task
	err := json.Unmarshal(response.Payload, &teetask)
	assert.NoError(t, err)
	assert.Equal(t, taskID, []byte(teetask.ID))

	response3 := stub.MockInvoke("3", [][]byte{[]byte(task.MethodGet), []byte("taskID1")})
	assert.Equal(t, shim.ERROR, int(response3.Status))
	assert.Contains(t, response3.Message, "Data does not exist")
}

func Test_getAll(t *testing.T) {
	stub := getTeeTaskMockStub()

	taskID1 := createTeeTask(t, "1", stub)
	taskID2 := createTeeTask(t, "2", stub)
	taskID3 := createTeeTask(t, "3", stub)

	response := stub.MockInvoke("4", [][]byte{[]byte(task.MethodGetAll)})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 3)

	var task1 tee.Task
	err1 := json.Unmarshal(dataList[0], &task1)
	assert.NoError(t, err1)
	assert.Equal(t, taskID1, []byte(task1.ID))

	var task2 tee.Task
	err2 := json.Unmarshal(dataList[1], &task2)
	assert.NoError(t, err2)
	assert.Equal(t, taskID2, []byte(task2.ID))

	var task3 tee.Task
	err3 := json.Unmarshal(dataList[2], &task3)
	assert.NoError(t, err3)
	assert.Equal(t, taskID3, []byte(task3.ID))
}

func getTeeTaskMockStub() *shim.MockStub {
	stub := shim.NewMockStub(task.ChaincodeName, new(TeeTaskChaincode))

	// Register tee chaincode
	teeStub := shim.NewMockStub(tee.ChaincodeName, new(MockTrustedExecutionEnv))
	stub.MockPeerChaincode(tee.ChaincodeName, teeStub)
	return stub
}

func getTeeTaskMockStubByTeeController(result []byte, errorMsg string) *shim.MockStub {
	stub := shim.NewMockStub(task.ChaincodeName, new(TeeTaskChaincode))

	teeStub := shim.NewMockStub(tee.ChaincodeName, &MockTEE{Result: result, ErrorMsg: errorMsg})
	stub.MockPeerChaincode(tee.ChaincodeName, teeStub)
	return stub
}
