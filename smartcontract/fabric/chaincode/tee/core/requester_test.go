package main

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
)

func Test_Invoke_request(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	filedEmptyError(t, stub, [][]byte{[]byte("request"), []byte(""), []byte("Requester")})
	filedEmptyError(t, stub, [][]byte{[]byte("request"), []byte("did"), []byte("")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("request"), []byte("Ciphertext"), []byte("Requester")})

	uploadData(t, stub, "1")

	response := requestData(t, stub, "3", "1")

	notification := notificationIsOk(t, response.Payload)

	bs, err := stub.GetState(notification.ID)
	assert.NoError(t, err)
	assert.Equal(t, bs, response.Payload)

	indexWorkWell(t, stub, requesterDIDIDIndex, []string{notification.Requester, notification.Data.ID})
	indexWorkWell(t, stub, requesterStatusDIDIDIndex, []string{notification.Requester, notification.Status.String(), notification.Data.ID})
	indexWorkWell(t, stub, ownerDIDIDIndex, []string{notification.Data.Owner, notification.Data.ID})
	indexWorkWell(t, stub, ownerRequesterDIDIDIndex, []string{notification.Data.Owner, notification.Requester, notification.Data.ID})
	indexWorkWell(t, stub, ownerStatusDIDIDIndex, []string{notification.Data.Owner, notification.Status.String(), notification.Data.ID})
	indexWorkWell(t, stub, ownerRequesterStatusDIDIDIndex, []string{notification.Data.Owner, notification.Requester, notification.Status.String(), notification.Data.ID})
}

func requestData(t *testing.T, stub *shim.MockStub, txID, did string) pb.Response {
	response := stub.MockInvoke(txID, [][]byte{[]byte("request"), []byte(did), []byte("Requester"), []byte(""), getEmptySigs(t)})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	preArgs := []string{did, PubHexForTest}
	hash, signatures := getHashAndSignatures(t, preArgs)
	response1 := stub.MockInvoke("requestID", [][]byte{[]byte("request"), []byte(did), []byte(PubHexForTest), []byte(hex.EncodeToString(hash)), signatures})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	return response
}

func request1Data(t *testing.T, stub *shim.MockStub, txID, did string) pb.Response {
	response := stub.MockInvoke(txID, [][]byte{[]byte("request"), []byte(did), []byte("Requester"), []byte(""), getEmptySigs(t)})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)
	return response
}

func notificationIsOk(t *testing.T, bs []byte) tee.Notification {
	var notification tee.Notification
	err := json.Unmarshal(bs, &notification)
	assert.NoError(t, err)

	dataIsOk(t, *notification.Data)
	assert.NotNil(t, notification.ID)
	assert.Equal(t, "Requester", notification.Requester)
	assert.NotZero(t, notification.RequestSecondsTimestamp)
	assert.Equal(t, tee.UnAuthorized, notification.Status)
	assert.Zero(t, notification.AuthSecondsTimestamp)
	return notification
}

func indexWorkWell(t *testing.T, stub *shim.MockStub, indexName string, attributes []string) {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, attributes)
	assert.NoError(t, err)
	defer resultsIterator.Close()
	assert.True(t, resultsIterator.HasNext())

	responseRange, err1 := resultsIterator.Next()
	assert.NoError(t, err1)

	objectType, compositeKeyParts, err2 := stub.SplitCompositeKey(responseRange.Key)
	assert.NoError(t, err2)
	assert.Equal(t, indexName, objectType)
	for index := 0; index < len(attributes); index++ {
		assert.Equal(t, attributes[index], compositeKeyParts[index])
	}
}

func Test_Invoke_queryRequestsByRequesterAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndDID"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndDID"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndDID"), []byte("Requester"), []byte("Ciphertext")})

	// query By Requester And DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryRequestsByRequesterAndDID"), []byte("Requester"), []byte("1")})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 2)

	notificationIsOk(t, dataList[0])
	notificationIsOk(t, dataList[1])

	// query By Requester Work Well
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryRequestsByRequesterAndDID"), []byte("Requester"), []byte("")})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 3)

	notificationIsOk(t, dataList1[0])
	notificationIsOk(t, dataList1[1])
	notificationIsOk(t, dataList1[2])
}

func Test_Invoke_queryRequestsByRequesterAndStatusAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte(""), []byte("0"), []byte("1")})
	filedEmptyError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Requester"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Ciphertext"), []byte("0"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Requester"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Requester"), []byte("0"), []byte("Ciphertext")})

	// query By Requester and status and DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Requester"), []byte("0"), []byte("1")})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 2)

	notificationIsOk(t, dataList[0])
	notificationIsOk(t, dataList[1])

	// query By Requester and status Work Well
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryRequestsByRequesterAndStatusAndDID"), []byte("Requester"), []byte("0"), []byte("")})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 3)

	notificationIsOk(t, dataList1[0])
	notificationIsOk(t, dataList1[1])
	notificationIsOk(t, dataList1[2])
}

func queryNotExistFieldError(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("2", args)
	assert.Equal(t, shim.ERROR, int(response.Status))
	assert.NotNil(t, response.Message)
}
