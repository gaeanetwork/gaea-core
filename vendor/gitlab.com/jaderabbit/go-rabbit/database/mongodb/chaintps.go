package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ChainTPSConnection for mongo db
type ChainTPSConnection struct {
	collection *mongo.Collection
}

var (
	connBChainTPS *ChainTPSConnection
	lockBChainTPS sync.Mutex
)

func GetChainTPSInfoConnection() (*ChainTPSConnection, error) {
	if connBChainTPS == nil {
		lockBChainTPS.Lock()
		defer lockBChainTPS.Unlock()

		if connBChainTPS == nil {
			Connect()
			if collection.ChainTPS == nil {
				return nil, errors.New("failed to connect tps collection")
			}

			connBChainTPS = &ChainTPSConnection{
				collection: collection.ChainTPS,
			}
		}
	}
	return connBChainTPS, nil
}

func (conn *ChainTPSConnection) Insert(chainTPS *types.ChainTPS) error {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	res, err := conn.collection.InsertOne(ctx, &chainTPS)
	if err != nil {
		return err
	}

	logger.Infof("insert tps collection successfully, %v", res.InsertedID)
	return nil
}

func (conn *ChainTPSConnection) Update(chainTPS *types.ChainTPS) error {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"_id": chainTPS.ID}, bson.M{"$set": &chainTPS}, op)
	if err != nil {
		return err
	}

	logger.Infof("update tps collection successfully, %v", res.UpsertedID)
	return nil
}

func (conn *ChainTPSConnection) InsertMany(chainTPSList []*types.ChainTPS) error {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	chainTPSs := []interface{}{}
	for _, chainTPS := range chainTPSList {
		chainTPSs = append(chainTPSs, chainTPS)
	}
	res, err := conn.collection.InsertMany(ctx, chainTPSs)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany chainTPSs collection successfully, %v", res.InsertedIDs)
	return nil
}

func (conn *ChainTPSConnection) Delete(id string) error {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	logger.Infof("delete chainTPS collection successfully, id:%s, DeletedCount:%d", id, res.DeletedCount)
	return nil
}

func (conn *ChainTPSConnection) DeleteMany(query interface{}) error {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany chainTPS collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *ChainTPSConnection) Count() (int64, error) {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the chainTPS collection successfully, %d", num)
	return num, nil
}

func (conn *ChainTPSConnection) QueryCount(query interface{}) (int64, error) {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the chainTPS collection successfully, %d", num)
	return num, nil
}

func (conn *ChainTPSConnection) Get(id string) (*types.ChainTPS, error) {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	chainTPS := &types.ChainTPS{}
	if err := conn.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(chainTPS); err != nil {
		return nil, err
	}
	return chainTPS, nil
}

func (conn *ChainTPSConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.ChainTPS, error) {
	lockBChainTPS.Lock()
	defer lockBChainTPS.Unlock()

	chainTPSs := []*types.ChainTPS{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"tps": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		chainTPS := &types.ChainTPS{}
		if err = cursor.Decode(chainTPS); err != nil {
			return nil, err
		}
		chainTPSs = append(chainTPSs, chainTPS)
	}
	return chainTPSs, nil
}
