package querier

import (
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// QueryHistoryByDID search data history with data id
func QueryHistoryByDID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error
	if err = chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	dataList := make([]*tee.SharedData, 0)
	for resultsIterator.HasNext() {
		var data tee.SharedData
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err = json.Unmarshal(response.Value, &data); err != nil {
			return shim.Error(err.Error())
		}

		dataList = append(dataList, &data)
	}
	if len(dataList) == 0 {
		return shim.Error(fmt.Sprintf("Data does not exist history, id: %v", id))
	}

	var bs []byte
	if bs, err = json.Marshal(dataList); err != nil {
		return shim.Error(fmt.Sprintf("Error marshaling data list: %v", err))
	}

	return shim.Success(bs)
}
