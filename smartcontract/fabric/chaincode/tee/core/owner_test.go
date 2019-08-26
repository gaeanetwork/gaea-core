package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jaderabbit/go-rabbit/tee"
)

// key pair for test
const (
	PrivHexForTest = "0x307702010104204368376222802d1a941f2eb0b7186a2c75f75e368946f923ad37e7c7718c2d7aa00a06082a8648ce3d030107a14403420004e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"
	PubHexForTest  = "0x04e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"
)

func Test_Invoke_queryNotificationsByOwnerAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndDID"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndDID"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndDID"), []byte("Owner"), []byte("Ciphertext")})

	// query By Owner and DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryNotificationsByOwnerAndDID"), []byte("Owner"), []byte("1")})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 3)

	notificationIsOk(t, dataList[0])
	notificationIsOk(t, dataList[2])

	// query By Owner Work Well
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryNotificationsByOwnerAndDID"), []byte("Owner"), []byte("")})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 5)

	notificationIsOk(t, dataList1[0])
	notificationIsOk(t, dataList1[2])
	notificationIsOk(t, dataList1[4])
}

func Test_Invoke_queryNotificationsByOwnerAndRequesterAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte(""), []byte("Requester"), []byte("1")})
	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Owner"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Ciphertext"), []byte("Requester"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Owner"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Owner"), []byte("Requester"), []byte("Ciphertext")})

	// query By Requester and status and DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Owner"), []byte("Requester"), []byte("1")})
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
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndDID"), []byte("Owner"), []byte("Requester"), []byte("")})
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

func Test_Invoke_queryNotificationsByOwnerAndStatusAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte(""), []byte("0"), []byte("1")})
	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Owner"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Ciphertext"), []byte("0"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Owner"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Owner"), []byte("0"), []byte("Ciphertext")})

	// query By Requester and status and DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Owner"), []byte("0"), []byte("1")})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err := json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 3)

	notificationIsOk(t, dataList[0])
	notificationIsOk(t, dataList[2])

	// query By Requester and status Work Well
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryNotificationsByOwnerAndStatusAndDID"), []byte("Owner"), []byte("0"), []byte("")})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)
	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 5)

	notificationIsOk(t, dataList1[0])
	notificationIsOk(t, dataList1[2])
	notificationIsOk(t, dataList1[4])
}

func Test_Invoke_queryNotificationsByOwnerAndRequesterAndStatusAndDID(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	requestData(t, stub, "3", "1")
	requestData(t, stub, "4", "2")
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte(""), []byte("Requester"), []byte("0"), []byte("1")})
	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte(""), []byte("0"), []byte("1")})
	filedEmptyError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Requester"), []byte(""), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Ciphertext"), []byte("Requester"), []byte("0"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Ciphertext"), []byte("0"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Requester"), []byte("Ciphertext"), []byte("1")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Requester"), []byte("0"), []byte("Ciphertext")})

	// query By Requester and status and DID Work Well
	response := stub.MockInvoke("6", [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Requester"), []byte("0"), []byte("1")})
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
	response1 := stub.MockInvoke("7", [][]byte{[]byte("queryNotificationsByOwnerAndRequesterAndStatusAndDID"), []byte("Owner"), []byte("Requester"), []byte("0"), []byte("")})
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

func Test_Invoke_authorize(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	response := requestData(t, stub, "3", "1")
	var data tee.Notification
	err := json.Unmarshal(response.Payload, &data)
	assert.NoError(t, err)
	response1 := requestData(t, stub, "4", "2")
	var data0 tee.Notification
	err0 := json.Unmarshal(response1.Payload, &data0)
	assert.NoError(t, err0)
	requestData(t, stub, "5", "1")

	filedEmptyError(t, stub, [][]byte{[]byte("authorize"), []byte(""), []byte("1"), []byte("Ciphertext"), []byte("0"), []byte("encryptedKey"), []byte("0")})
	filedEmptyError(t, stub, [][]byte{[]byte("authorize"), []byte("3"), []byte(""), []byte("Ciphertext"), []byte("0"), []byte("encryptedKey"), []byte("0")})
	filedEmptyError(t, stub, [][]byte{[]byte("authorize"), []byte("3"), []byte("1"), []byte(""), []byte("0"), []byte("encryptedKey"), []byte("0")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("authorize"), []byte("Ciphertext"), []byte("1"), []byte("Ciphertext"), []byte("encryptedKey"), []byte("0")})
	queryNotExistFieldError(t, stub, [][]byte{[]byte("authorize"), []byte("3"), []byte("Ciphertext"), []byte("Ciphertext"), []byte("encryptedKey"), []byte("0")})

	authorizeData(t, stub, "6", data.ID, "1")
	bs1, err1 := stub.GetState(data.ID)
	assert.NoError(t, err1)
	assert.NotEmpty(t, bs1)

	var data1 tee.Notification
	err11 := json.Unmarshal(bs1, &data1)
	assert.NoError(t, err11)
	assert.NotNil(t, data1.Data)
	assert.Equal(t, "Ciphertext", data1.DataInfo.DataStoreAddress)
	assert.Equal(t, "1", data1.Status.String())
	assert.NotZero(t, data1.AuthSecondsTimestamp)
	assert.Empty(t, data1.RefusedReason)

	authorizeData(t, stub, "7", data0.ID, "2")
	bs2, err2 := stub.GetState(data0.ID)
	assert.NoError(t, err2)
	assert.NotEmpty(t, bs2)

	var data2 tee.Notification
	err21 := json.Unmarshal(bs2, &data2)
	assert.NoError(t, err21)
	assert.NotNil(t, data2.Data)
	assert.Equal(t, "Ciphertext", data2.RefusedReason)
	assert.Equal(t, "2", data2.Status.String())
	assert.NotZero(t, data2.AuthSecondsTimestamp)
	assert.Empty(t, data2.DataInfo)
}

func authorize1Data(t *testing.T, stub *shim.MockStub, txID, id, status string) {
	args := [][]byte{[]byte("authorize"), []byte(id), []byte(status), []byte("Ciphertext"), []byte("0"), []byte("encryptedKey"), []byte("0"), []byte(""), getEmptySigs(t)}
	response := stub.MockInvoke(txID, args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.Nil(t, response.Payload)
}

func authorizeData(t *testing.T, stub *shim.MockStub, txID, id, status string) {
	args := [][]byte{[]byte("authorize"), []byte(id), []byte(status), []byte("Ciphertext"), []byte("0"), []byte("encryptedKey"), []byte("0"), []byte(""), getEmptySigs(t)}
	response := stub.MockInvoke(txID, args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.Nil(t, response.Payload)

	idData := append([]byte("requestID"), []byte("1")...)
	idHash := sha256.Sum256(idData)
	args1 := [][]byte{[]byte("authorize"), []byte(hex.EncodeToString(idHash[:])), []byte(status), []byte("Ciphertext"), []byte("0"), []byte("encryptedKey"), []byte("0")}
	response1 := stub.MockInvoke("authorizeID", args1)
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.Nil(t, response1.Payload)
}

func Test_Invoke_updateAuthIndex(t *testing.T) {
	stub := shim.NewMockStub("tee", new(TrustedExecutionEnv))

	uploadData(t, stub, "1")
	uploadData(t, stub, "2")
	response := request1Data(t, stub, "3", "1")
	var data tee.Notification
	err := json.Unmarshal(response.Payload, &data)
	assert.NoError(t, err)
	response1 := request1Data(t, stub, "4", "2")
	var data0 tee.Notification
	err0 := json.Unmarshal(response1.Payload, &data0)
	assert.NoError(t, err0)
	request1Data(t, stub, "5", "1")

	authorize1Data(t, stub, "6", data.ID, "1")
	authorize1Data(t, stub, "7", data0.ID, "2")

	// request status 0
	response = queryRequestsByRequesterAndStatusAndDID(stub, []string{"Requester", "0", ""})
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Payload)

	var dataList [][]byte
	err = json.Unmarshal(response.Payload, &dataList)
	assert.NoError(t, err)
	assert.Len(t, dataList, 1)
	notificationIsOk(t, dataList[0])

	// request status 1
	response1 = queryRequestsByRequesterAndStatusAndDID(stub, []string{"Requester", "1", ""})
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
	assert.NotNil(t, response1.Payload)

	var dataList1 [][]byte
	err1 := json.Unmarshal(response1.Payload, &dataList1)
	assert.NoError(t, err1)
	assert.Len(t, dataList1, 1)

	var notification tee.Notification
	err2 := json.Unmarshal(dataList1[0], &notification)
	assert.NoError(t, err2)

	dataIsOk(t, *notification.Data)
	assert.Equal(t, data.ID, notification.ID)
	assert.Equal(t, tee.Authorized, notification.Status)
	assert.NotZero(t, notification.AuthSecondsTimestamp)
	assert.Equal(t, "Ciphertext", notification.DataInfo.DataStoreAddress)

	// request status 2
	response2 := queryRequestsByRequesterAndStatusAndDID(stub, []string{"Requester", "2", ""})
	assert.Equal(t, shim.OK, int(response2.Status))
	assert.Empty(t, response2.Message)
	assert.NotNil(t, response2.Payload)

	var dataList2 [][]byte
	err3 := json.Unmarshal(response2.Payload, &dataList2)
	assert.NoError(t, err3)
	assert.Len(t, dataList2, 1)

	var notification1 tee.Notification
	err4 := json.Unmarshal(dataList2[0], &notification1)
	assert.NoError(t, err4)

	dataIsOk(t, *notification1.Data)
	assert.Equal(t, data0.ID, notification1.ID)
	assert.Equal(t, tee.Refused, notification1.Status)
	assert.NotZero(t, notification1.AuthSecondsTimestamp)
	assert.Equal(t, "Ciphertext", notification1.RefusedReason)

	// other index
	response3 := queryNotificationsByOwnerAndStatusAndDID(stub, []string{"Owner", "0", ""})
	assert.Equal(t, response, response3)
	response4 := queryNotificationsByOwnerAndStatusAndDID(stub, []string{"Owner", "1", ""})
	assert.Equal(t, response1, response4)
	response5 := queryNotificationsByOwnerAndStatusAndDID(stub, []string{"Owner", "2", ""})
	assert.Equal(t, response2, response5)

	response6 := queryNotificationsByOwnerAndRequesterAndStatusAndDID(stub, []string{"Owner", "Requester", "0", ""})
	assert.Equal(t, response, response6)
	response7 := queryNotificationsByOwnerAndRequesterAndStatusAndDID(stub, []string{"Owner", "Requester", "1", ""})
	assert.Equal(t, response1, response7)
	response8 := queryNotificationsByOwnerAndRequesterAndStatusAndDID(stub, []string{"Owner", "Requester", "2", ""})
	assert.Equal(t, response2, response8)
}
