package querier

import (
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// QueryDataByID search all data with id
func QueryDataByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
