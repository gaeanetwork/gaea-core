package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TxInfoConnection for mongo db
type TxInfoConnection struct {
	collection *mongo.Collection
}

var (
	connTx *TxInfoConnection
	lockTx sync.Mutex
)

func GetTxInfoConnection() (*TxInfoConnection, error) {
	if connTx == nil {
		lockTx.Lock()
		defer lockTx.Unlock()
		if connTx == nil {
			Connect()
			if collection.TxInfo == nil {
				return nil, errors.New("failed to connect txinfo collection")
			}

			connTx = &TxInfoConnection{
				collection: collection.TxInfo,
			}
		}
	}
	return connTx, nil
}

func (conn *TxInfoConnection) Insert(tx *types.TxInfo) error {
	lockTx.Lock()
	defer lockTx.Unlock()

	res, err := conn.collection.InsertOne(ctx, &tx)
	if err != nil {
		return err
	}

	logger.Infof("insert txinfo collection successfully, %v", res.InsertedID)
	return nil
}

func (conn *TxInfoConnection) update(tx *types.TxInfo) error {
	lockTx.Lock()
	defer lockTx.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"_id": tx.ID}, bson.M{"$set": &tx}, op)
	if err != nil {
		return err
	}

	logger.Infof("update txinfo collection successfully, %v", res.UpsertedID)
	return nil
}

func (conn *TxInfoConnection) InsertMany(txList []*types.TxInfo) error {
	lockTx.Lock()
	defer lockTx.Unlock()

	txs := []interface{}{}
	for _, tx := range txList {
		txs = append(txs, tx)
	}
	res, err := conn.collection.InsertMany(ctx, txs)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany txinfo collection successfully, %v", res.InsertedIDs)
	return nil
}

func (conn *TxInfoConnection) Delete(txID string) error {
	lockTx.Lock()
	defer lockTx.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.M{"txid": txID})
	if err != nil {
		return err
	}

	logger.Infof("delete txinfo collection successfully, id:%s, DeletedCount:%d", txID, res.DeletedCount)
	return nil
}

func (conn *TxInfoConnection) DeleteMany(query interface{}) error {
	lockTx.Lock()
	defer lockTx.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany txinfo collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *TxInfoConnection) Count() (int64, error) {
	lockTx.Lock()
	defer lockTx.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the txinfo collection successfully, %d", num)
	return num, nil
}

func (conn *TxInfoConnection) QueryCount(query interface{}) (int64, error) {
	lockTx.Lock()
	defer lockTx.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the txinfo collection successfully, %d", num)
	return num, nil
}

func (conn *TxInfoConnection) Get(txID string) (*types.TxInfo, error) {
	lockTx.Lock()
	defer lockTx.Unlock()

	tx := &types.TxInfo{}
	if err := conn.collection.FindOne(ctx, bson.M{"txid": txID}).Decode(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (conn *TxInfoConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.TxInfo, error) {
	lockTx.Lock()
	defer lockTx.Unlock()

	txs := []*types.TxInfo{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"timestamp": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		tx := &types.TxInfo{}
		if err = cursor.Decode(tx); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (conn *TxInfoConnection) Aggregate(pipeline interface{}) ([]*types.UserGroupTx, error) {
	lockTx.Lock()
	defer lockTx.Unlock()

	ugts := []*types.UserGroupTx{}
	cursor, err := conn.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		ugt := &types.UserGroupTx{}
		if err = cursor.Decode(ugt); err != nil {
			return nil, err
		}
		ugts = append(ugts, ugt)
	}
	return ugts, nil
}
