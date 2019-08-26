package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// index value
var (
	CompositeValue = []byte{0x00}
)

type getDataByID func(stub shim.ChaincodeStubInterface, args []string) peer.Response

// CheckArgsEmpty check chaincode args are non-empty
func CheckArgsEmpty(args []string, length int) error {
	if l := len(args); l < length {
		return fmt.Errorf("Incorrect number of arguments. Expecting be greater than or equal to %d, Actual: %d(%v)", length, l, args)
	}

	for index := 0; index < length; index++ {
		if len(args[index]) <= 0 {
			return fmt.Errorf("The index %d argument must be a non-empty string, args: %v", index, args)
		}
	}
	return nil
}

// CheckArgsLength check the length of chaincode args
func CheckArgsLength(args []string, length int) error {
	if l := len(args); l < length {
		return fmt.Errorf("Incorrect number of arguments. Expecting be greater than or equal to %d, Actual: %d(%v)", length, l, args)
	}

	return nil
}

// GetDataByID query data by id
func GetDataByID(stub shim.ChaincodeStubInterface, id string) pb.Response {
	var bs []byte
	var err error
	if bs, err = stub.GetState(id); err != nil {
		return shim.Error("Failed to get state for id: " + id)
	} else if len(bs) == 0 {
		return shim.Error("Data does not exist, id: " + id)
	}

	return shim.Success(bs)
}

// GetDataListByIndexAndKeys through index to query data list
func GetDataListByIndexAndKeys(stub shim.ChaincodeStubInterface, index string, keys []string, idIndex int, getFunc getDataByID) pb.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get state by partial composite key, index: %s, keys: %v, error: %v", index, keys, err))
	}
	defer resultsIterator.Close()

	count, dataList := 0, make([][]byte, 0)
	var response pb.Response
	for ; resultsIterator.HasNext(); count++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to split composite key: %s", responseRange.Key))
		}

		if response = getFunc(stub, []string{compositeKeyParts[idIndex]}); response.Status != shim.OK {
			return shim.Error(fmt.Sprintf("Failed to query data by id: %s", response.Message))
		}

		dataList = append(dataList, response.Payload)
	}
	if count == 0 {
		return shim.Error(fmt.Sprintf("Data does not exist, index: %s, keys: %v", index, keys))
	}

	var bs []byte
	if bs, err = json.Marshal(dataList); err != nil {
		return shim.Error(fmt.Sprintf("Error marshaling data list: %v", err))
	}
	return shim.Success(bs)
}
