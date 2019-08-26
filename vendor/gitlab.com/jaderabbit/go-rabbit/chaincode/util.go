package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common"
)

type getDataByID func(stub shim.ChaincodeStubInterface, args []string) pb.Response

// CCMethod the method in the chaincode
type CCMethod func(stub shim.ChaincodeStubInterface, args []string) pb.Response

// index value
var (
	CompositeValue = []byte{0x00}
)

// GetCCMethod get the chaincode method
func GetCCMethod(ccMapMethod map[string]CCMethod, methodName string) (CCMethod, error) {
	if len(methodName) == 0 {
		return nil, errors.New("not specified methodName")
	}

	if ccMapMethod == nil {
		return nil, errors.New("not specified ccMapMethod")
	}

	ccMethod, ok := ccMapMethod[methodName]
	if ok {
		return ccMethod, nil
	}

	var builder strings.Builder
	for key := range ccMapMethod {
		builder.WriteString(" `")
		builder.WriteString(key)
		builder.WriteString("` ")
	}

	return nil, fmt.Errorf("Invalid invoke function name:%s, Expecting:%s", methodName, builder.String())
}

func CreateComposite(stub shim.ChaincodeStubInterface, key string, attributes []string) error {
	indexKey, err := stub.CreateCompositeKey(key, attributes)
	if err != nil {
		return err
	}

	value := []byte{0x00}
	return stub.PutState(indexKey, value)
}

func DeleteComposite(stub shim.ChaincodeStubInterface, key string, attributes []string) error {
	indexKey, err := stub.CreateCompositeKey(key, attributes)
	if err != nil {
		return err
	}

	return stub.DelState(indexKey)
}

func PutState(stub shim.ChaincodeStubInterface, key string, obj interface{}) ([]byte, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err:%s", err.Error())
	}

	if err := stub.PutState(key, objBytes); err != nil {
		return nil, fmt.Errorf("stub.PutState error, key:%s, err:%s", key, err.Error())
	}
	return objBytes, nil
}

func GetState(stub shim.ChaincodeStubInterface, key string, obj interface{}) error {
	objByte, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("stub.GetState key:%s, err:%s", key, err.Error())
	}

	if len(objByte) == 0 {
		return fmt.Errorf("key(%s) does not exist", key)
	}

	if err = json.Unmarshal(objByte, obj); err != nil {
		return fmt.Errorf("json.Unmarshal key:%s, err:%s", key, err.Error())
	}
	return nil
}

// ConstructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
func ConstructQueryResponseFromIterator(stub shim.ChaincodeStubInterface, resultsIterator shim.StateQueryIteratorInterface, getIndex int) ([]string, error) {
	arrayData := []string{}
	if resultsIterator == nil {
		return arrayData, nil
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, err
		}

		arrayData = append(arrayData, compositeKeyParts[getIndex])
	}

	arrayData = common.ReverseArray(arrayData)

	return arrayData, nil
}

// GetMSPID returns the ID associated with the invoking identity.  This ID
// is guaranteed to be unique within the MSP.
func GetMSPID(stub shim.ChaincodeStubInterface) (string, error) {
	mspid, err := cid.GetID(stub)
	if err != nil {
		if strings.Contains(err.Error(), "identity bytes are neither X509 PEM format nor an idemix credential") {
			return "msp", nil
		}
	}

	return mspid, err
}

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
