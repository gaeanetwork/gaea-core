package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// queryNotificationsByOwnerAndDID search all owner notifications, if did! = "", returns the owner to specify the data id notifications
func queryNotificationsByOwnerAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The owner argument must be a non-empty string")
	}

	owner, did := args[0], args[1]
	var keys []string
	if did == "" {
		keys = []string{owner}
	} else {
		keys = []string{owner, did}
	}

	return queryDataListByIndexAndKeys(stub, ownerDIDIDIndex, keys, 2)
}

// queryNotificationsByOwnerAndRequesterAndDID search all owner to specify the requester notifications, if did! = "",
// notifications that returns the specified data ID of the specified requester of the owner
func queryNotificationsByOwnerAndRequesterAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The owner argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("The requester argument must be a non-empty string")
	}

	owner, requester, did := args[0], args[1], args[2]
	var keys []string
	if did == "" {
		keys = []string{owner, requester}
	} else {
		keys = []string{owner, requester, did}
	}

	return queryDataListByIndexAndKeys(stub, ownerRequesterDIDIDIndex, keys, 3)
}

// queryNotificationsByOwnerAndStatusAndDID search all owner to specify the status notifications, if did! = "",
// notifications that returns the specified data ID of the specified status of the owner
func queryNotificationsByOwnerAndStatusAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The owner argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("The status argument must be a non-empty string")
	}

	owner, status, did := args[0], args[1], args[2]
	var keys []string
	if did == "" {
		keys = []string{owner, status}
	} else {
		keys = []string{owner, status, did}
	}

	return queryDataListByIndexAndKeys(stub, ownerStatusDIDIDIndex, keys, 3)
}

// queryNotificationsByOwnerAndRequesterAndStatusAndDID search all owner to specify the requester and status notifications, if did! = "",
// notifications that returns the specified data ID of the specified status of the specified requester of the owner
func queryNotificationsByOwnerAndRequesterAndStatusAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The owner argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("The requester argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("The status argument must be a non-empty string")
	}

	owner, requester, status, did := args[0], args[1], args[2], args[3]
	var keys []string
	if did == "" {
		keys = []string{owner, requester, status}
	} else {
		keys = []string{owner, requester, status, did}
	}

	return queryDataListByIndexAndKeys(stub, ownerRequesterStatusDIDIDIndex, keys, 4)
}

func authorize(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 3); err != nil {
		return shim.Error(fmt.Sprintf("Error checking Args[id, status, message], error: %v", err))
	}

	id, status, message := args[0], args[1], args[2]
	response := queryDataByID(stub, []string{id})
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}

	var oldNotification tee.Notification
	if err = json.Unmarshal(response.Payload, &oldNotification); err != nil {
		return shim.Error("Error unmarshaling notification: " + err.Error())
	}

	// If the upload is signed, the update and authorize must also be signed.
	if data := oldNotification.Data; len(data.Signatures) > 0 && data.Signatures[0] != "" {
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, data.Owner); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	newNotification := oldNotification
	switch status {
	case "1":
		newNotification.Status = tee.Authorized

		if err = chaincode.CheckArgsEmpty(args, 6); err != nil {
			return shim.Error(fmt.Sprintf("Error checking Args[id, status, message, dataStoreTypeStr, encryptedKey, encryptedTypeStr], error: %v", err))
		}
		dataStoreTypeStr, encryptedKey, encryptedTypeStr := args[3], args[4], args[5]

		dataStoreType, err := strconv.Atoi(dataStoreTypeStr)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to parse data store type to int. dataStoreType: %s, error: %v", dataStoreTypeStr, err))
		}

		encryptedType, err := strconv.Atoi(encryptedTypeStr)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to parse encrypted type to int. encryptedType: %s, error: %v", encryptedTypeStr, err))
		}

		newNotification.DataInfo = &tee.DataInfo{
			DataStoreAddress: message,
			DataStoreType:    tee.DataStoreType(dataStoreType),
			EncryptedKey:     encryptedKey,
			EncryptedType:    tee.EncryptedType(encryptedType),
		}
	case "2":
		newNotification.Status = tee.Refused
		newNotification.RefusedReason = message
	default:
		return shim.Error("Invalid authorize status. Expecting \"1\": Authorized \"2\": Refused. Actual: " + status)
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}
	newNotification.AuthSecondsTimestamp = timestamp.Seconds

	var bs []byte
	if bs, err = json.Marshal(newNotification); err != nil {
		return shim.Error("Error marshalingnotification: " + err.Error())
	}

	if err = stub.PutState(newNotification.ID, bs); err != nil {
		return shim.Error("Error putting shard data to state: " + err.Error())
	}

	return updateAuthIndexes(stub, oldNotification, newNotification)
}

func updateAuthIndexes(stub shim.ChaincodeStubInterface, oldNotification, newNotification tee.Notification) pb.Response {
	var err error
	// Save the index from the requester's perspective  [requester~status~did~id]
	indexKey, attributes := "", []string{newNotification.Requester, newNotification.Status.String(), newNotification.Data.ID, newNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(requesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~status~did~id]
	indexKey, attributes = "", []string{newNotification.Data.Owner, newNotification.Status.String(), newNotification.Data.ID, newNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~requester~status~did~id]
	indexKey, attributes = "", []string{newNotification.Data.Owner, newNotification.Requester, newNotification.Status.String(), newNotification.Data.ID, newNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerRequesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// delete the index from the requester's perspective  [requester~status~did~id]
	indexKey, attributes = "", []string{oldNotification.Requester, oldNotification.Status.String(), oldNotification.Data.ID, oldNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(requesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.DelState(indexKey)

	// Save the index from the owner's perspective  [owner~status~did~id]
	indexKey, attributes = "", []string{oldNotification.Data.Owner, oldNotification.Status.String(), oldNotification.Data.ID, oldNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.DelState(indexKey)

	// Save the index from the owner's perspective  [owner~requester~status~did~id]
	indexKey, attributes = "", []string{oldNotification.Data.Owner, oldNotification.Requester, oldNotification.Status.String(), oldNotification.Data.ID, oldNotification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerRequesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.DelState(indexKey)

	return shim.Success(nil)
}
