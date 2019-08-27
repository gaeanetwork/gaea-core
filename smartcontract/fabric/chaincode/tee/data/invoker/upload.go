package invoker

import (
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/data/data"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// upload shared sharedData
func upload(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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

	sharedData := &tee.SharedData{
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
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, sharedData.Owner); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	var bs []byte
	if bs, err = json.Marshal(sharedData); err != nil {
		return shim.Error("Error marshaling shard sharedData: " + err.Error())
	}

	if err = stub.PutState(sharedData.ID, bs); err != nil {
		return shim.Error("Error putting shard sharedData to state: " + err.Error())
	}

	// Save ownerIDIndexKey
	var ownerIDIndexKey string
	if ownerIDIndexKey, err = stub.CreateCompositeKey(data.OwnerIDIndex, []string{sharedData.Owner, sharedData.ID}); err != nil {
		return shim.Error(err.Error())
	}

	stub.PutState(ownerIDIndexKey, data.CompositeValue)
	return shim.Success(bs)
}
