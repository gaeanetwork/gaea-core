package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ActionLogConnection for mongo db
type ActionLogConnection struct {
	collection *mongo.Collection
}

var (
	connActionLog *ActionLogConnection
	lockActionLog sync.Mutex
)

func GetActionLogConnection() (*ActionLogConnection, error) {
	if connActionLog == nil {
		lockActionLog.Lock()
		defer lockActionLog.Unlock()

		if connActionLog == nil {
			Connect()
			if collection.ActionLog == nil {
				return nil, errors.New("failed to connect ActionLog collection")
			}

			connActionLog = &ActionLogConnection{
				collection: collection.ActionLog,
			}
		}
	}
	return connActionLog, nil
}

func (conn *ActionLogConnection) Insert(actionLog *types.ActionLog) error {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	res, err := conn.collection.InsertOne(ctx, &actionLog)
	if err != nil {
		return err
	}

	logger.Infof("insert actionLog collection successfully, %v", res.InsertedID)
	return nil
}

func (conn *ActionLogConnection) InsertMany(actionLogList []*types.ActionLog) error {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	actionLogs := []interface{}{}
	for _, actionLog := range actionLogList {
		actionLogs = append(actionLogs, actionLog)
	}
	res, err := conn.collection.InsertMany(ctx, actionLogs)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany actionLogs collection successfully, %v", res.InsertedIDs)
	return nil
}

func (conn *ActionLogConnection) Delete(id string) error {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	logger.Infof("delete actionLog collection successfully, id:%s, DeletedCount:%d", id, res.DeletedCount)
	return nil
}

func (conn *ActionLogConnection) DeleteMany(query interface{}) error {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany actionLog collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *ActionLogConnection) Count() (int64, error) {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the actionLog collection successfully, %d", num)
	return num, nil
}

func (conn *ActionLogConnection) QueryCount(query interface{}) (int64, error) {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the actionLog collection successfully, %d", num)
	return num, nil
}

func (conn *ActionLogConnection) Get(id string) (*types.ActionLog, error) {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	actionLog := &types.ActionLog{}
	if err := conn.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(actionLog); err != nil {
		return nil, err
	}
	return actionLog, nil
}

func (conn *ActionLogConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.ActionLog, error) {
	lockActionLog.Lock()
	defer lockActionLog.Unlock()

	actionLogs := []*types.ActionLog{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"timestamp": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		actionLog := &types.ActionLog{}
		if err = cursor.Decode(actionLog); err != nil {
			return nil, err
		}
		actionLogs = append(actionLogs, actionLog)
	}
	return actionLogs, nil
}
