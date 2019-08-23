package task

import (
	"fmt"
	"log"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee/task"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// ExecutionService chaincode tee task implementation
type ExecutionService struct {
}

// Init the system chaincode, save the system key pair
func (t *ExecutionService) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	var err error
	if err = chaincode.CheckArgsEmpty(args, 2); err != nil {
		return shim.Error(err.Error())
	}

	privHex, pubHex := args[0], args[1]
	stub.PutState(task.KeyPrivHex, []byte(privHex))
	stub.PutState(task.KeyPubHex, []byte(pubHex))
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *ExecutionService) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	go func() { log.Println(fmt.Sprintf("TxID: %s, function: %s, args: %v", stub.GetTxID(), function, args)) }()
	switch function {
	case task.MethodCreate:
	case task.MethodGet:
	case task.MethodGetAll:
	case task.MethodExecute:
	}

	return shim.Error("Invalid invoke function name. Expecting \"create\" \"get\" \"getAll\" \"execute\", Actual: " + function)
}
