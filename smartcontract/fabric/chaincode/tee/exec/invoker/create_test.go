package invoker

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/container"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

func Test_create(t *testing.T) {
	stub := getTeeTaskMockStub()

	// checkIfAlgorithmAndDataExistInTEE Parameters
	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodCreate)})
	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodCreate), []byte("")})
	IncorrectNumberArgsError(t, stub, [][]byte{[]byte(task.MethodCreate), []byte(""), []byte("")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodCreate), []byte("dataIDs"), []byte(""), []byte("")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodCreate), []byte(""), []byte("algorithmID"), []byte("")})
	filedEmptyError(t, stub, [][]byte{[]byte(task.MethodCreate), []byte(""), []byte(""), []byte("resultAddress")})

	// Empty DataIDs
	response1 := stub.MockInvoke("1", [][]byte{[]byte(task.MethodCreate), []byte("ids"), []byte("algorithmID"), []byte("resultAddress")})
	assert.Equal(t, shim.ERROR, int(response1.Status))
	assert.Contains(t, response1.Message, "Failed to parse dataIDs to string slice")

	dataIDs := []string{}
	ids, err := json.Marshal(dataIDs)
	assert.NoError(t, err)
	response2 := stub.MockInvoke("2", [][]byte{[]byte(task.MethodCreate), ids, []byte("algorithmID"), []byte("/tmp/teetask/test")})
	assert.Equal(t, shim.ERROR, int(response2.Status))
	assert.Contains(t, response2.Message, "Error task data IDs must be non-empty")

	// Right Input
	dataIDs1 := []string{"id1", "id2"}
	ids1, err1 := json.Marshal(dataIDs1)
	assert.NoError(t, err1)
	response3 := stub.MockInvoke("3", [][]byte{[]byte(task.MethodCreate), ids1, []byte("algorithmID"), []byte("/tmp/teetask/test")})
	assert.Equal(t, shim.OK, int(response3.Status))
	assert.Empty(t, response3.Message)

	// checkIfAlgorithmAndDataExistInTEE task
	response4 := stub.MockInvoke("4", [][]byte{[]byte("get"), response3.Payload})
	assert.Equal(t, shim.OK, int(response4.Status))
	assert.Empty(t, response4.Message)

	var task tee.Task
	err = json.Unmarshal(response4.Payload, &task)
	assert.NoError(t, err)
	assert.Equal(t, dataIDs1, task.DataIDs)
	assert.Equal(t, "algorithmID", task.AlgorithmID)
	assert.NotZero(t, task.CreateSecondsTimestamp)
	assert.NotZero(t, task.UploadSecondsTimestamp)

}

func Test_Check(t *testing.T) {
	stub, teetask := getTeeTaskMockStub(), getTeeTaskForTests()

	teetask.DataNotifications = nil
	err := checkIfAlgorithmAndDataExistInTEE(stub, teetask)
	assert.NoError(t, err)

	// DataIDs is not exist and query err
	teetask = getTeeTaskForTests()
	teetask.DataIDs[0] = "err"
	err = checkIfAlgorithmAndDataExistInTEE(stub, teetask)
	assert.Contains(t, err.Error(), "error: err")

	teetask = getTeeTaskForTests()
	teetask.DataIDs[0] = "jsonerr"
	err = checkIfAlgorithmAndDataExistInTEE(stub, teetask)
	assert.Contains(t, err.Error(), "invalid character")

	// AlgorithmID is not exist and query err
	teetask = getTeeTaskForTests()
	teetask.AlgorithmID = "err"
	err = checkIfAlgorithmAndDataExistInTEE(stub, teetask)
	assert.Contains(t, err.Error(), "error: err")

	teetask = getTeeTaskForTests()
	teetask.AlgorithmID = "jsonerr"
	err8 := checkIfAlgorithmAndDataExistInTEE(stub, teetask)
	assert.Contains(t, err8.Error(), "invalid character")
}

func Test_RequestData(t *testing.T) {
	stub, teetask := getTeeTaskMockStub(), getTeeTaskForTests()

	err1 := requestAuthorizationsForAllData(stub, teetask)
	assert.NoError(t, err1)

	// DataIDs is not exist and query err
	teetask1 := getTeeTaskForTests()
	teetask1.DataIDs[0] = "err"
	err4 := requestAuthorizationsForAllData(stub, teetask1)
	assert.Contains(t, err4.Error(), "error: err")

	teetask2 := getTeeTaskForTests()
	teetask2.DataIDs[0] = "jsonerr"
	err5 := requestAuthorizationsForAllData(stub, teetask2)
	assert.Contains(t, err5.Error(), "invalid character")
}

func getTeeTaskForTests() *tee.Task {
	return &tee.Task{
		ID:                "1",
		DataIDs:           []string{"id1", "id2"},
		AlgorithmID:       "algorithmID",
		Container:         container.Docker,
		ResultAddress:     "resultAddress",
		DataNotifications: make(map[string]string),
		Partners:          make(map[string]struct{}),
	}
}

func Test_RequestDataWithSignature(t *testing.T) {
	stub, teetask := getTeeTaskMockStub(), getTeeTaskForTests()

	teetask.DataIDs[0] = "sigID1"
	// Init privHex error
	err := requestAuthorizationsForAllData(stub, teetask)
	assert.Contains(t, err.Error(), "Error getting byte slice from privHex")

	// Success invoke
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})
	err1 := requestAuthorizationsForAllData(stub, teetask)
	assert.NoError(t, err1)

	// DataIDs is not exist and query err
	teetask1 := getTeeTaskForTests()
	teetask1.DataIDs[0] = "err"
	err2 := requestAuthorizationsForAllData(stub, teetask1)
	assert.Contains(t, err2.Error(), "error: err")

	teetask2 := getTeeTaskForTests()
	teetask2.DataIDs[0] = "jsonerr"
	err3 := requestAuthorizationsForAllData(stub, teetask2)
	assert.Contains(t, err3.Error(), "invalid character")
}

func createTeeTask(t *testing.T, uuid string, stub *shim.MockStub) []byte {
	dataIDs := []string{"id1", "id2"}
	ids, err := json.Marshal(dataIDs)
	assert.NoError(t, err)
	response := stub.MockInvoke(uuid, [][]byte{[]byte(task.MethodCreate), ids, []byte("algorithmID"), []byte("/tmp/teetask/test")})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	return response.Payload
}
