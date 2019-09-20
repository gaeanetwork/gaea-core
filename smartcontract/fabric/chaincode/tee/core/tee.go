package main

import (
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func main() {
	err := shim.Start(new(TrustedExecutionEnv))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

var (
	ownerIDIndex                   = "owner~id"
	requesterDIDIDIndex            = "requester~did~id"
	requesterStatusDIDIDIndex      = "requester~status~did~id"
	ownerDIDIDIndex                = "owner~did~id"
	ownerRequesterDIDIDIndex       = "owner~requester~did~id"
	ownerStatusDIDIDIndex          = "owner~status~did~id"
	ownerRequesterStatusDIDIDIndex = "owner~requester~status~did~id"
	compositeValue                 = []byte{0x00}
)

// TrustedExecutionEnv chaincode implementation
type TrustedExecutionEnv struct {
}

// Init the system chaincode
func (t *TrustedExecutionEnv) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *TrustedExecutionEnv) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case tee.MethodUpload:
		return upload(stub, args)
	case tee.MethodUpdate:
		return update(stub, args)
	case tee.MethodQueryDataByID:
		return queryDataByID(stub, args)
	case tee.MethodQueryHistoryByDID:
		return queryHistoryByDID(stub, args)
	case tee.MethodQueryDataByOwner:
		return queryDataByOwner(stub, args)
	case tee.MethodRequest:
		return request(stub, args)
	case tee.MethodQueryRequestsByRequesterAndDID:
		return queryRequestsByRequesterAndDID(stub, args)
	case tee.MethodQueryRequestsByRequesterAndStatusAndDID:
		return queryRequestsByRequesterAndStatusAndDID(stub, args)
	case tee.MethodQueryNotificationsByOwnerAndDID:
		return queryNotificationsByOwnerAndDID(stub, args)
	case tee.MethodQueryNotificationsByOwnerAndRequesterAndDID:
		return queryNotificationsByOwnerAndRequesterAndDID(stub, args)
	case tee.MethodQueryNotificationsByOwnerAndStatusAndDID:
		return queryNotificationsByOwnerAndStatusAndDID(stub, args)
	case tee.MethodQueryNotificationsByOwnerAndRequesterAndStatusAndDID:
		return queryNotificationsByOwnerAndRequesterAndStatusAndDID(stub, args)
	case tee.MethodAuthorize:
		return authorize(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"upload\" \"update\" \"queryDataByID\" \"queryDataByOwner\" \"request\" \"queryRequestsByRequesterAndDID\" \"queryRequestsByRequesterAndStatusAndDID\" \"queryNotificationsByOwnerAndDID\" \"queryNotificationsByOwnerAndRequesterAndDID\" \"queryNotificationsByOwnerAndStatusAndDID\" \"queryNotificationsByOwnerAndRequesterAndStatusAndDID\" Actual: " + function)
}

// upload shared data
func upload(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 4); err != nil {
		return shim.Error(err.Error())
	}
	var sigs []string
	if len(args) > 5 {
		if err = json.Unmarshal([]byte(args[5]), &sigs); err != nil {
			return shim.Error("Failed to unmarshal signatures for args[5], error: " + err.Error())
		} else if len(sigs) == 0 {
			return shim.Error("Signature slice length is 0")
		}
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}

	data := &tee.SharedData{
		ID:                     stub.GetTxID(),
		Ciphertext:             args[0],
		Hash:                   args[1],
		Description:            args[2],
		Owner:                  args[3],
		Signatures:             sigs,
		CreateSecondsTimestamp: timestamp.Seconds,
		UpdateSecondsTimestamp: timestamp.Seconds,
	}

	if len(sigs) > 0 && sigs[0] != "" {
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, data.Owner); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	var bs []byte
	if bs, err = json.Marshal(data); err != nil {
		return shim.Error("Error marshaling shard data: " + err.Error())
	}

	if err = stub.PutState(data.ID, bs); err != nil {
		return shim.Error("Error putting shard data to state: " + err.Error())
	}

	// Save ownerIDIndexKey
	var ownerIDIndexKey string
	if ownerIDIndexKey, err = stub.CreateCompositeKey(ownerIDIndex, []string{data.Owner, data.ID}); err != nil {
		return shim.Error(err.Error())
	}

	stub.PutState(ownerIDIndexKey, compositeValue)
	return shim.Success(bs)
}

// update ciphertext, hash, description of shared data
func update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 4); err != nil {
		return shim.Error(err.Error())
	}
	var sigs []string
	if len(args) > 5 {
		if err = json.Unmarshal([]byte(args[5]), &sigs); err != nil {
			return shim.Error("Failed to unmarshal signatures for args[5], error: " + err.Error())
		} else if len(sigs) == 0 {
			return shim.Error("Signature slice length is 0")
		}
	}

	id, ciphertext, hash, description, bs := args[0], args[1], args[2], args[3], []byte{}
	if bs, err = stub.GetState(id); err != nil {
		return shim.Error("Failed to get state for id: " + id)
	} else if len(bs) == 0 {
		return shim.Error("Shared data does not exist, id: " + id)
	}

	var data tee.SharedData
	if err = json.Unmarshal(bs, &data); err != nil {
		return shim.Error("Error unmarshaling shard data: " + err.Error())
	}

	// If the upload is signed, the update and authorize must also be signed.
	if len(data.Signatures) > 0 && data.Signatures[0] != "" {
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, data.Owner); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}

		data.Signatures = sigs
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}

	data.Ciphertext = ciphertext
	data.Hash = hash
	data.Description = description
	data.UpdateSecondsTimestamp = timestamp.Seconds
	if bs, err = json.Marshal(data); err != nil {
		return shim.Error("Error marshaling shard data: " + err.Error())
	}

	if err = stub.PutState(data.ID, bs); err != nil {
		return shim.Error("Error putting shard data to state: " + err.Error())
	}

	return shim.Success(nil)
}

// queryDataByID search all data with id
func queryDataByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	id, bs := args[0], []byte{}
	if bs, err = stub.GetState(id); err != nil {
		return shim.Error("Failed to get state for id: " + id)
	} else if len(bs) == 0 {
		return shim.Error("Shared data does not exist, id: " + id)
	}

	return shim.Success(bs)
}

// queryHistoryByDID search data history with data id
func queryHistoryByDID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	dataList := make([]*tee.SharedData, 0)
	for resultsIterator.HasNext() {
		var data tee.SharedData
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err = json.Unmarshal(response.Value, &data); err != nil {
			return shim.Error(err.Error())
		}

		dataList = append(dataList, &data)
	}
	if len(dataList) == 0 {
		return shim.Error(fmt.Sprintf("Data does not exist history, id: %v", id))
	}

	var bs []byte
	if bs, err = json.Marshal(dataList); err != nil {
		return shim.Error(fmt.Sprintf("Error marshaling data list: %v", err))
	}

	return shim.Success(bs)
}

// queryDataByOwner search all data with owner
func queryDataByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	owner := args[0]
	return queryDataListByIndexAndKeys(stub, ownerIDIndex, []string{owner}, 1)
}

func queryDataListByIndexAndKeys(stub shim.ChaincodeStubInterface, index string, keys []string, idIndex int) pb.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get state by partial composite key, index: %s, keys: %v, error: %v", index, keys, err))
	}
	defer resultsIterator.Close()

	count, dataList := 0, make([][]byte, 0)
	var response pb.Response
	for ; resultsIterator.HasNext(); count++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to split composite key: %s", responseRange.Key))
		}

		if response = queryDataByID(stub, []string{compositeKeyParts[idIndex]}); response.Status != shim.OK {
			return shim.Error(fmt.Sprintf("Failed to query data by id: %s", response.Message))
		}

		dataList = append(dataList, response.Payload)
	}
	if count == 0 {
		return shim.Error(fmt.Sprintf("Data does not exist, index: %s, keys: %v", index, keys))
	}

	var bs []byte
	if bs, err = json.Marshal(dataList); err != nil {
		return shim.Error(fmt.Sprintf("Error marshaling data list: %v", err))
	}
	return shim.Success(bs)
}
