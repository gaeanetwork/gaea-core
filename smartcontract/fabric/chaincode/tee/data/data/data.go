package data

import (
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/data/querier"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// index keys
const (
	OwnerIDIndex                   = "owner~id"
	RequesterDIDIDIndex            = "requester~did~id"
	RequesterStatusDIDIDIndex      = "requester~status~did~id"
	OwnerDIDIDIndex                = "owner~did~id"
	OwnerRequesterDIDIDIndex       = "owner~requester~did~id"
	OwnerStatusDIDIDIndex          = "owner~status~did~id"
	OwnerRequesterStatusDIDIDIndex = "owner~requester~status~did~id"
)

// CompositeValue index value
var CompositeValue = []byte{0x00}

// SharedDataService chaincode implementation
type SharedDataService struct {
}

// Init the system chaincode
func (t *SharedDataService) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *SharedDataService) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case tee.MethodUpload:
	case tee.MethodUpdate:
	case tee.MethodQueryDataByID:
		return querier.QueryDataByID(stub, args)
	case tee.MethodQueryHistoryByDID:
		return querier.QueryHistoryByDID(stub, args)
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
