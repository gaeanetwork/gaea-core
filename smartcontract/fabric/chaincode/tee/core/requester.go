package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// request request shared data
func request(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 2); err != nil {
		return shim.Error(err.Error())
	}

	did, requester := args[0], args[1]
	response := queryDataByID(stub, []string{did})
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}

	var data tee.SharedData
	if err = json.Unmarshal(response.Payload, &data); err != nil {
		return shim.Error("Error unmarshaling shard data: " + err.Error())
	}

	// If the upload is signed, the update, authorize and request must also be signed.
	if len(data.Signatures) > 0 && data.Signatures[0] != "" {
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, requester); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}

	// The reason we don't use the tx ID directly is that task is to use all data IDs requested
	// by the same transaction context. This means that the tx IDs of all requests are the same.
	idData := append([]byte(stub.GetTxID()), []byte(data.ID)...)
	idHash := sha256.Sum256(idData)
	notification := &tee.Notification{
		ID:                      hex.EncodeToString(idHash[:]),
		Data:                    &data,
		Requester:               requester,
		RequestSecondsTimestamp: timestamp.Seconds,
		Status:                  tee.UnAuthorized,
	}

	var bs []byte
	if bs, err = json.Marshal(notification); err != nil {
		return shim.Error("Error marshaling notification: " + err.Error())
	}

	if err = stub.PutState(notification.ID, bs); err != nil {
		return shim.Error("Error putting notification to state: " + err.Error())
	}

	// Save the index from the requester's perspective  [requester~did~id]
	indexKey, attributes := "", []string{notification.Requester, notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(requesterDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the requester's perspective  [requester~status~did~id]
	indexKey, attributes = "", []string{notification.Requester, notification.Status.String(), notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(requesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~did~id]
	indexKey, attributes = "", []string{notification.Data.Owner, notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~requester~did~id]
	indexKey, attributes = "", []string{notification.Data.Owner, notification.Requester, notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerRequesterDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~status~did~id]
	indexKey, attributes = "", []string{notification.Data.Owner, notification.Status.String(), notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	// Save the index from the owner's perspective  [owner~requester~status~did~id]
	indexKey, attributes = "", []string{notification.Data.Owner, notification.Requester, notification.Status.String(), notification.Data.ID, notification.ID}
	if indexKey, err = stub.CreateCompositeKey(ownerRequesterStatusDIDIDIndex, attributes); err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(indexKey, compositeValue)

	return shim.Success(bs)
}

// queryRequestsByRequesterAndDID search all requester requests, if did! = "", returns the requester to specify the data id requests
func queryRequestsByRequesterAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The requester argument must be a non-empty string")
	}

	requester, did := args[0], args[1]
	var keys []string
	if did == "" {
		keys = []string{requester}
	} else {
		keys = []string{requester, did}
	}

	return queryDataListByIndexAndKeys(stub, requesterDIDIDIndex, keys, 2)
}

// queryRequestsByRequesterAndStatusAndDID search all requester status requests, if did! = "", returns the requester to specify the data id requests
func queryRequestsByRequesterAndStatusAndDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if len(args[0]) <= 0 {
		return shim.Error("The requester argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("The status argument must be a non-empty string")
	}

	requester, status, did := args[0], args[1], args[2]
	var keys []string
	if did == "" {
		keys = []string{requester, status}
	} else {
		keys = []string{requester, status, did}
	}

	return queryDataListByIndexAndKeys(stub, requesterStatusDIDIDIndex, keys, 3)
}
