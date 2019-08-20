package querier

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func queryDataListByIndexAndKeys(stub shim.ChaincodeStubInterface, index string, keys []string, idIndex int) peer.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get state by partial composite key, index: %s, keys: %v, error: %v", index, keys, err))
	}
	defer resultsIterator.Close()

	count, dataList := 0, make([][]byte, 0)
	var response peer.Response
	for ; resultsIterator.HasNext(); count++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to split composite key: %s", responseRange.Key))
		}

		if response = QueryDataByID(stub, []string{compositeKeyParts[idIndex]}); response.Status != shim.OK {
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
