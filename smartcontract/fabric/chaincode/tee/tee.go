package main

import (
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
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
	case tee.MethodUpdate:
	case tee.MethodQueryDataByID:
		return queryDataByID(stub, args)
	case tee.MethodQueryHistoryByDID:
	case tee.MethodQueryDataByOwner:
	case tee.MethodRequest:
	case tee.MethodQueryRequestsByRequesterAndDID:
	case tee.MethodQueryRequestsByRequesterAndStatusAndDID:
	case tee.MethodQueryNotificationsByOwnerAndDID:
	case tee.MethodQueryNotificationsByOwnerAndRequesterAndDID:
	case tee.MethodQueryNotificationsByOwnerAndStatusAndDID:
	case tee.MethodQueryNotificationsByOwnerAndRequesterAndStatusAndDID:
	case tee.MethodAuthorize:
	}

	return shim.Error("Invalid invoke function name. Expecting \"upload\" \"update\" \"queryDataByID\" \"queryDataByOwner\" \"request\" \"queryRequestsByRequesterAndDID\" \"queryRequestsByRequesterAndStatusAndDID\" \"queryNotificationsByOwnerAndDID\" \"queryNotificationsByOwnerAndRequesterAndDID\" \"queryNotificationsByOwnerAndStatusAndDID\" \"queryNotificationsByOwnerAndRequesterAndStatusAndDID\" Actual: " + function)
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
