package src

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/asset"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/sdk"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/system"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto"
	"gitlab.com/jaderabbit/go-rabbit/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	syncLock sync.Mutex

	aesKey = []byte("yi$l1an~r-tong!_")
	desKey = []byte("yilianrongtong2019011613")
	sm4Key = []byte("yi$lian^r*tong@_")

	logger = flogging.MustGetLogger("src.asset")
)

// SyncToMongoDB sync asst to mongodb which invokes successfully
func SyncToMongoDB(key string) {
	assetDB, err := mongodb.GetAssetConnection()
	if err != nil {
		logger.Errorf("failed to get asset mongodb connection, err:%s", err.Error())
		return
	}

	asset, err := QueryAsset(key)
	if err != nil {
		logger.Errorf("failed to get asset, err:%s, assetKey:%s", err.Error(), key)
		return
	}

	asset.MapData = nil

	if err = assetDB.Update(asset); err != nil {
		logger.Errorf("failed to update data to MongoDB, key:%s, err:%s", key, err.Error())
		return
	}
	logger.Debugf("successful synchronization of asset:%s", key)
}

// QueryAsset Query the asset by key
func QueryAsset(key string, strMUName ...string) (*asset.Asset, error) {
	mspUserName := sdk.DefaultMspUserName
	if len(strMUName) > 0 && len(strMUName[0]) > 0 {
		mspUserName = strMUName[0]
	}

	chaincodeServer, err := sdk.GetAssetServer(mspUserName)
	if err != nil {
		return nil, err
	}

	arrayStr := []string{"query", key}
	str, err := chaincodeServer.Query(arrayStr)
	if err != nil {
		return nil, err
	}

	if len(str) == 0 {
		return nil, fmt.Errorf("asset does not exist, key:%s", key)
	}

	newAsset := &asset.Asset{}
	if err = json.Unmarshal([]byte(str), newAsset); err != nil {
		return nil, err
	}
	return newAsset, nil
}

// QueryAssetList query the assets by filter, eg:key, txID, userID, beginTime, endTime
func QueryAssetList(pageIndex, pageSize uint32, key, txID, userID string, beginTime, endTime int64) (*system.QueryResult, error) {
	skip := (pageIndex - 1) * pageSize

	mapM := make(bson.M)
	if len(key) > 0 {
		mapM["key"] = key
	}

	if len(txID) > 0 {
		mapM["tx_id"] = txID
	}

	if len(userID) > 0 {
		mapM["operate_userid"] = userID
	}

	if beginTime > 0 && endTime > 0 {
		mapM["timestamp"] = bson.M{"$gte": beginTime, "$lte": endTime}
	} else if beginTime > 0 {
		mapM["timestamp"] = bson.M{"$lte": endTime}
	} else if endTime > 0 {
		mapM["timestamp"] = bson.M{"$lte": endTime}
	}

	assetDB, err := mongodb.GetAssetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get asset mongodb connection, err:%s", err.Error())
	}

	count, err := assetDB.QueryCount(mapM)
	if err != nil {
		return nil, err
	}

	arrayAsset, err := assetDB.GetByFilter(mapM, int64(skip), int64(pageSize))
	if err != nil {
		return nil, err
	}

	result := &system.QueryResult{
		TotalNumber: int(count),
		Data:        arrayAsset,
		PageIndex:   int(pageIndex),
		PageSize:    int(pageSize),
	}
	return result, nil
}

// SM2DataVerification sm2 encrypt data verify
func SM2DataVerification(originData, dstData, encryptKey string) (bool, error) {
	key, err := common.HexToBytes(encryptKey)
	if err != nil {
		return false, err
	}

	byteExistDstData, err := base64.StdEncoding.DecodeString(dstData)
	if err != nil {
		return false, err
	}

	ciphertext, err := crypto.SM2Decrypt(byteExistDstData, key)
	if err != nil {
		return false, err
	}

	if originData == string(ciphertext) {
		return true, nil
	}

	return false, err
}

// Dncrypt decrypt data by key and data
func Dncrypt(encryptHash crypto.Hash, data, key []byte) ([]byte, error) {
	handler, err := crypto.GetDecrypt(encryptHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get the decrypt method, err:%s, encryptHash:%d", err.Error(), encryptHash)
	}
	return handler(data, key)
}
