package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/chaincode/asset"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//AssetConnection for mongo db
type AssetConnection struct {
	collection *mongo.Collection
}

var (
	connAsset *AssetConnection
	lockAsset sync.Mutex
)

// GetAssetConnection get the asset connection from the mongodb
func GetAssetConnection() (*AssetConnection, error) {
	if connAsset == nil {
		lockAsset.Lock()
		defer lockAsset.Unlock()
		if connAsset == nil {
			Connect()
			if collection.Asset == nil {
				return nil, errors.New("failed to connect asset collection")
			}

			connAsset = &AssetConnection{
				collection: collection.Asset,
			}
		}
	}

	return connAsset, nil
}

// Insert insert the asset document
func (conn *AssetConnection) Insert(asset *asset.Asset) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	res, err := conn.collection.InsertOne(ctx, &asset)
	if err != nil {
		return err
	}

	logger.Infof("insert asset collection successfully, %v", res.InsertedID)
	return nil
}

// Update insert the asset document
func (conn *AssetConnection) Update(asset *asset.Asset) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"key": asset.Key}, bson.M{"$set": &asset}, op)
	if err != nil {
		return err
	}

	logger.Infof("update asset collection successfully, %v", res.UpsertedID)
	return nil
}

// UpdateMany update the asset document
func (conn *AssetConnection) UpdateMany(assets []*asset.Asset) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	arrayBSONM := []bson.M{}
	for _, at := range assets {
		arrayBSONM = append(arrayBSONM, bson.M{"key": at.Key})
	}

	filter := bson.M{"$or": arrayBSONM}
	res, err := conn.collection.UpdateMany(ctx, filter, bson.M{"$set": assets}, op)
	if err != nil {
		return err
	}

	logger.Infof("UpdateMany asset collection successfully, %v", res.UpsertedID)
	return nil
}

// InsertMany insert the asset documents
func (conn *AssetConnection) InsertMany(assetList []*asset.Asset) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	assets := []interface{}{}
	for _, asset := range assetList {
		assets = append(assets, asset)
	}
	res, err := conn.collection.InsertMany(ctx, assets)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany asset collection successfully, %v", res.InsertedIDs)
	return nil
}

// Delete delete the asset ducument by id
func (conn *AssetConnection) Delete(id string) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.M{"key": id})
	if err != nil {
		return err
	}

	logger.Infof("delete asset collection successfully, id:%s, DeletedCount:%d", id, res.DeletedCount)
	return nil
}

// DeleteMany delete the asset documents by query filter
func (conn *AssetConnection) DeleteMany(query interface{}) error {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany asset collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

// Count query the number of all asset document in the mongodb
func (conn *AssetConnection) Count() (int64, error) {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the asset collection successfully, %d", num)
	return num, nil
}

// QueryCount query the number of the asset document by query filter in the mongodb
func (conn *AssetConnection) QueryCount(query interface{}) (int64, error) {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Debugf("query count of the asset collection successfully, %d", num)
	return num, nil
}

// Get get the asset by id
func (conn *AssetConnection) Get(id string) (*asset.Asset, error) {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	asset := &asset.Asset{}
	if err := conn.collection.FindOne(ctx, bson.M{"key": id}).Decode(asset); err != nil {
		return nil, err
	}
	return asset, nil
}

// GetByFilter get the asset by filter
func (conn *AssetConnection) GetByFilter(query interface{}, skip, limit int64) ([]*asset.Asset, error) {
	lockAsset.Lock()
	defer lockAsset.Unlock()

	assets := []*asset.Asset{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"create_time": 1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		asset := &asset.Asset{}
		if err = cursor.Decode(asset); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}
