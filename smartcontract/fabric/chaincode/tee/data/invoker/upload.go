package invoker

import (
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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
		UploadSecondsTimestamp: timestamp.Seconds,
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
