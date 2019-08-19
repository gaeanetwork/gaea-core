package data

import (
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/data/querier"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
