package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

var (
	logger = flogging.MustGetLogger("TeeTaskChaincode")
)

func main() {
	err := shim.Start(new(TeeTaskChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// TeeTaskChaincode chaincode tee task implementation
type TeeTaskChaincode struct {
}

// Init the system chaincode, save the system key pair
func (t *TeeTaskChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
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
func (t *TeeTaskChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	go func() { logger.Infof("TxID: %s, function: %s, args: %v", stub.GetTxID(), function, args) }()
	switch function {
	case task.MethodCreate:
		return create(stub, args)
	case task.MethodGet:
		return get(stub, args)
	case task.MethodGetAll:
		return getAll(stub, args)
	case task.MethodExecute:
		return execute(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"create\" \"get\" \"getAll\" \"execute\", Actual: " + function)
}

func get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	dataID := args[0]
	if dataID == task.KeyPrivHex {
		return shim.Error(errors.Wrapf(err, "Invalid data id: %s", dataID).Error())
	}

	return chaincode.GetDataByID(stub, dataID)
}

func getAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return chaincode.GetDataListByIndexAndKeys(stub, task.TaskIDIndex, []string{task.ChaincodeName}, 1, get)
}

func saveTaskAndIndex(stub shim.ChaincodeStubInterface, teetask *tee.Task, saveIndex bool) error {
	taskBytes, err := json.Marshal(teetask)
	if err != nil {
		return errors.Wrap(err, "failed to parse tee task to bytes")
	}

	if err = stub.PutState(teetask.ID, taskBytes); err != nil {
		return errors.Wrap(err, "failed to put tee task to state")
	}

	if saveIndex {
		taskIDIndexKey, err := stub.CreateCompositeKey(task.TaskIDIndex, []string{task.ChaincodeName, teetask.ID})
		if err != nil {
			return errors.Wrap(err, "failed to create composite key to save tee task index")
		}

		return stub.PutState(taskIDIndexKey, chaincode.CompositeValue)
	}

	return nil
}
