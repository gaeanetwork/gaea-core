package algorithm

import (
	"github.com/gaeanetwork/gaea-core/protos/tee"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

// upload tee algorithm
func upload(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	var timestamp *timestamp.Timestamp
	if timestamp, err = stub.GetTxTimestamp(); err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}

	var algo tee.Algorithm
	if err = proto.Unmarshal([]byte(args[0]), &algo); err != nil {
		return shim.Error("Invalid args[0], cannot unmarshal to algo: " + err.Error())
	}

	algo.Id = stub.GetTxID()
	algo.CreateSeconds = timestamp.Seconds
	algo.UpdateSeconds = timestamp.Seconds

	if len(algo.Signatures) > 0 && algo.Signatures[0] != "" {
		// TODO - CHECK CHANGE TO VERIFY HASH
		if err = chaincode.CheckArgsContainsHashAndSignatures(args, algo.Owner); err != nil {
			return shim.Error("Failed to verify ecdsa signature, error: " + err.Error())
		}
	}

	var bs []byte
	if bs, err = proto.Marshal(&algo); err != nil {
		return shim.Error("Error marshaling algo: " + err.Error())
	}

	// Save data
	if err = stub.PutState(algo.Id, bs); err != nil {
		return shim.Error("Error putting algo to state: " + err.Error())
	}

	// Save Indexs
	saveIndexFn := func() error {
		allIDKey, err := stub.CreateCompositeKey(IndexAllID, []string{ChaincodeName, algo.Id})
		if err != nil {
			return errors.Wrapf(err, "failed to create IndexAllID compositeKey, id: %s", algo.Id)
		}
		stub.PutState(allIDKey, chaincode.CompositeValue)

		nameIDKey, err := stub.CreateCompositeKey(IndexNameID, []string{algo.Name, algo.Id})
		if err != nil {
			return errors.Wrapf(err, "failed to create IndexNameID compositeKey, id: %s, name: %s", algo.Id, algo.Name)
		}
		stub.PutState(nameIDKey, chaincode.CompositeValue)

		return nil
	}

	if err = saveIndexFn(); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(algo.Id))
}
