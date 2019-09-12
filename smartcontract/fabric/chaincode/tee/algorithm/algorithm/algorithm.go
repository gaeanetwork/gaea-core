package algorithm

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// index keys
const (
	// ChaincodeName
	ChaincodeName = "tee_algo"

	// Index keys
	IndexAllID  = "algorithm~id"
	IndexNameID = "name~id"

	// Methods
	MethodUpload               = "uploadAlgorithm"
	MethodQueryAlgorithmByName = "queryAlgorithmByName"
	MethodQueryAlgorithms      = "queryAllAlgorithms"
)

// CompositeValue index value
var CompositeValue = []byte{0x00}

// ChaincodeService chaincode implementation
type ChaincodeService struct {
}

// Init the system chaincode
func (t *ChaincodeService) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *ChaincodeService) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case MethodUpload:
		return upload(stub, args)
	case MethodQueryAlgorithmByName:
	case MethodQueryAlgorithms:
	}

	return shim.Error("Invalid invoke function name. Expecting \"upload\" \"queryAlgorithmByName\" \"queryAllAlgorithms\" Actual: " + function)
}
