package asset

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
)

var (
	// CCMapMethod a method set about user chaincode
	CCMapMethod = map[string]chaincode.CCMethod{
		"add":               add,
		"update":            update,
		"verify":            verify,
		"query":             query,
		"history":           getHistory,
		"managerfieldright": managerFieldRight,
		"manager":           manager,
		"queryright":        queryRightInfo,
		"submitbatch":       SubmitBatch,
		"querybatch":        queryBatch,
	}
)

// Chaincode the financial chain of supplier
type Chaincode struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// inadvertently clobber your ledger’s data!
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	mspID, err := chaincode.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get msp id, err:%s", err.Error()))
	}

	if err = stub.PutState(theOwnerMspIDofAsset, []byte(mspID)); err != nil {
		return shim.Error(fmt.Sprintf("failed to put msp id to db, err:%s", err.Error()))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	ccMethod, err := chaincode.GetCCMethod(CCMapMethod, function)
	if err != nil {
		return shim.Error(err.Error())
	}

	return ccMethod(stub, args)
}

// add the data
// mapData format: {"field1":{"data":"=sdsd=wweasdf=","encrypt_type":4},"field2":{"data":"8a8sdfaf","encrypt_type":5}}
// field1 is the fieldname, =sdsd=wweasdf= is the ciphertext
// encrypt_type is the encrypt type
func add(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 3); err != nil {
		return shim.Error(err.Error())
	}

	key := args[0]

	mapData := make(map[string]*Field)
	if err := json.Unmarshal([]byte(args[1]), &mapData); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal data, err:%s, data:%s", err.Error(), args[1]))
	}

	public, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to parse public field, public(arg[2]):%s", args[2]))
	}

	assetByte, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get asset, key:%s, err:%s", key, err.Error()))
	}

	if len(assetByte) > 0 {
		return shim.Error(fmt.Sprintf("failed to add asset, the key(%s) has already exist", key))
	}

	newAsset := NewAsset()

	user, err := getUserByChaincode(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get user by chaincode, err:%s", err.Error()))
	}

	newAsset.Key = key
	newAsset.TxID = stub.GetTxID()
	newAsset.UserID = user.ID
	newAsset.MapData = mapData
	newAsset.IsPublic = public

	assetByte, err = json.Marshal(newAsset)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to marshal asset, err:%s", err.Error()))
	}

	// put the asset to chain
	if err := stub.PutState(key, assetByte); err != nil {
		return shim.Error(fmt.Sprintf("failed to put state, key:%s, err:%s", key, err.Error()))
	}
	return shim.Success(nil)
}

// SubmitBatch the data
// mapData format: {"field1":{"data":"=sdsd=wweasdf=","encrypt_type":4},"field2":{"data":"8a8sdfaf","encrypt_type":5}}
// field1 is the fieldname, =sdsd=wweasdf= is the ciphertext
// encrypt_type is the encrypt type
func SubmitBatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	var arrayDatas []*Asset
	if err := json.Unmarshal([]byte(args[0]), &arrayDatas); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal data, err:%s", err.Error()))
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get user by chaincode, err:%s", err.Error()))
	}

	// check the key and data, remove the repeat asset
	mapAssetRemoveRepeat := make(map[string]*Asset)
	for _, newAsset := range arrayDatas {
		if len(newAsset.Key) == 0 {
			continue
		}

		if len(newAsset.MapData) == 0 {
			continue
		}

		_, ok := mapAssetRemoveRepeat[newAsset.Key]
		if ok {
			continue
		}
		mapAssetRemoveRepeat[newAsset.Key] = newAsset
	}

	for _, newAsset := range mapAssetRemoveRepeat {
		newAsset.TxID = stub.GetTxID()
		newAsset.UserID = user.ID

		assetByte, err := json.Marshal(newAsset)
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to marshal asset, key:%s, err:%s", newAsset.Key, err.Error()))

		}

		// put the asset to chain
		if err := stub.PutState(newAsset.Key, assetByte); err != nil {
			return shim.Error(fmt.Sprintf("failed to put state, key:%s, err:%s", newAsset.Key, err.Error()))
		}
	}

	// wg := sync.WaitGroup{}

	// // get the error info
	// errorInfoChan := make(chan string, 100)
	// errInfoArray := []string{}
	// go func() {
	// 	for {
	// 		select {
	// 		case errorInfo := <-errorInfoChan:
	// 			errInfoArray = append(errInfoArray, errorInfo)
	// 		default:
	// 		}
	// 	}
	// }()

	// for _, newAsset := range mapAssetRemoveRepeat {
	// 	wg.Add(1)

	// 	go func(at *Asset) {
	// 		defer wg.Done()
	// 		at.TxID = stub.GetTxID()
	// 		at.UserID = user.ID

	// 		assetByte, err := json.Marshal(at)
	// 		if err != nil {
	// 			errorInfoChan <- fmt.Sprintf("failed to marshal asset, key:%s, err:%s", at.Key, err.Error())
	// 			return
	// 		}

	// 		// put the asset to chain
	// 		if err := stub.PutState(at.Key, assetByte); err != nil {
	// 			errorInfoChan <- fmt.Sprintf("failed to put state, key:%s, err:%s", at.Key, err.Error())
	// 		}
	// 	}(newAsset)
	// }
	// wg.Wait()

	// if len(errInfoArray) > 0 {
	// 	byteErrInfo, err := json.Marshal(&errInfoArray)
	// 	if err != nil {
	// 		shim.Error(fmt.Sprintf("failed to json marshal errInfoArray, InfoArray:%v, err:%s", errInfoArray, err.Error()))
	// 	}
	// 	return shim.Error(string(byteErrInfo))
	// }
	return shim.Success(nil)
}

func update(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsLength(args, 4); err != nil {
		return shim.Error(err.Error())
	}

	key := args[0]

	mapData := make(map[string]*Field)
	if len(args[1]) > 0 {
		if err := json.Unmarshal([]byte(args[1]), &mapData); err != nil {
			return shim.Error(fmt.Sprintf("failed to unmarshal data, err:%s, data:%s", err.Error(), args[1]))
		}
	}

	removeFields := []string{}
	if len(args[2]) > 0 {
		if err := json.Unmarshal([]byte(args[2]), &removeFields); err != nil {
			return shim.Error(fmt.Sprintf("failed to unmarshal remove fields, err:%s, removeFields:%s", err.Error(), args[2]))
		}
	}

	var public bool
	if len(args[3]) > 0 {
		var err error
		public, err = strconv.ParseBool(args[3])
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to parse public field, public(arg[3]):%s", args[3]))
		}
	}

	assetByte, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get asset, key:%s, err:%s", key, err.Error()))
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get user by chaincode, err:%s", err.Error()))
	}

	newAsset := NewAsset()
	err = json.Unmarshal(assetByte, newAsset)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal asset, err:%s, str:%s", err.Error(), string(assetByte)))
	}

	if newAsset.UserID != user.ID {
		return shim.Error(fmt.Sprintf("No permission to update the data, key:%s", key))
	}

	newAsset.TxID = stub.GetTxID()

	for fieldName, field := range mapData {
		newAsset.MapData[fieldName] = field
	}

	if len(removeFields) > 0 {
		for _, fieldName := range removeFields {
			delete(newAsset.MapData, fieldName)
		}
	}

	if len(args[3]) > 0 {
		newAsset.IsPublic = public
	}

	assetByte, err = json.Marshal(newAsset)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to marshal asset, err:%s", err.Error()))
	}

	// put the asset to chain
	if err := stub.PutState(key, assetByte); err != nil {
		return shim.Error(fmt.Sprintf("failed to put state, key:%s, err:%s", key, err.Error()))
	}
	return shim.Success(nil)
}

func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}
	key := args[0]

	assetByte, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("stub.GetState(key) key:%s, err:%s", key, err.Error()))
	}

	if len(assetByte) == 0 {
		return shim.Success(nil)
	}

	newAsset := &Asset{}
	if err = json.Unmarshal(assetByte, newAsset); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal asset, err:%s, str:%s", err.Error(), string(assetByte)))
	}

	if err = newAsset.filter(stub); err != nil {
		return shim.Error(fmt.Sprintf("failed to filter asset, err:%s", err.Error()))
	}

	assetByte, err = json.Marshal(newAsset)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to marshal asset, err:%s", err.Error()))
	}
	return shim.Success(assetByte)
}

func queryBatch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	var arrayKey []string
	if err := json.Unmarshal([]byte(args[0]), &arrayKey); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal keys to array, keys:%s, err:%s", args[0], err.Error()))
	}

	mapKeysRemoveRepeat := make(map[string]struct{})
	for _, key := range arrayKey {
		if len(key) == 0 {
			continue
		}

		if _, ok := mapKeysRemoveRepeat[key]; ok {
			continue
		}
		mapKeysRemoveRepeat[key] = struct{}{}
	}

	var arrayDatas []*Asset
	for key := range mapKeysRemoveRepeat {
		assetByte, err := stub.GetState(key)
		if err != nil {
			return shim.Error(fmt.Sprintf("stub.GetState(key) key:%s, err:%s", key, err.Error()))
		}

		if len(assetByte) == 0 {
			continue
		}

		newAsset := &Asset{}
		if err = json.Unmarshal(assetByte, newAsset); err != nil {
			return shim.Error(fmt.Sprintf("failed to unmarshal asset, err:%s, str:%s", err.Error(), string(assetByte)))
		}

		if err = newAsset.filter(stub); err != nil {
			return shim.Error(fmt.Sprintf("failed to filter asset, err:%s", err.Error()))
		}
		arrayDatas = append(arrayDatas, newAsset)
	}

	assetByte, err := json.Marshal(arrayDatas)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to marshal asset, err:%s", err.Error()))
	}
	return shim.Success(assetByte)
}

func verify(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 2); err != nil {
		return shim.Error(err.Error())
	}

	key := args[0]

	mapData := make(map[string]*Field)
	if len(args[1]) > 0 {
		if err := json.Unmarshal([]byte(args[1]), &mapData); err != nil {
			return shim.Error(fmt.Sprintf("failed to ummarshal mapData(arg[1]), err:%s", err.Error()))
		}
	}

	// get the asset from chain by key
	assetByte, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get state, key:%s", key))
	}

	if len(assetByte) == 0 {
		return shim.Error(fmt.Sprintf("key(%s) does not exist", key))
	}

	newAsset := &Asset{}
	if err = json.Unmarshal(assetByte, newAsset); err != nil {
		return shim.Error(fmt.Sprintf("failed to unmarshal asset, err:%s, str:%s", err.Error(), string(assetByte)))
	}

	for fieldName, field := range mapData {
		f, ok := newAsset.MapData[fieldName]
		if !ok {
			return shim.Error(fmt.Sprintf("failed to verify data, not existed field name: %s", fieldName))
		}

		if ok, err := f.Equal(field); err != nil {
			return shim.Error(err.Error())
		} else if !ok {
			return shim.Success([]byte(fmt.Sprint(false)))
		}
	}

	return shim.Success([]byte(fmt.Sprint(true)))
}

// getHistory get the page size history, default is 100
func getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	key := args[0]

	defaultPageSize := 100
	deepHistory := defaultPageSize
	if len(args) > 1 && len(args[1]) > 0 {
		var err error
		deepHistory, err = strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to convert arg[1] to int, err:%v, str:%s", err, args[1]))
		}
	}

	// get a history of key values across time
	iterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iterator.Close()

	arrayAsset := []*Asset{}

	skip := 0

	if deepHistory > defaultPageSize {
		skip = (deepHistory / defaultPageSize) * defaultPageSize
	}
	count := 0
	for iterator.HasNext() {
		count++
		if skip > count {
			continue
		}

		response, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		historyAsset := &Asset{}
		if err = json.Unmarshal(response.Value, historyAsset); err != nil {
			return shim.Error(err.Error())
		}

		if err = historyAsset.filter(stub); err != nil {
			return shim.Error(err.Error())
		}
		arrayAsset = append(arrayAsset, historyAsset)

		if count >= skip+defaultPageSize || count >= deepHistory {
			break
		}
	}

	historyByte, err := json.Marshal(&arrayAsset)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(historyByte)
}
