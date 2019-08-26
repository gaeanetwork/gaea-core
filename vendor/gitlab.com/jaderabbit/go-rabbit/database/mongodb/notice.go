package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//NoticeConnection for mongo db
type NoticeConnection struct {
	collection *mongo.Collection
}

var (
	connNotice *NoticeConnection
	lockNotice sync.Mutex
)

// GetNoticeConnection get the notice connection from the mongodb
func GetNoticeConnection() (*NoticeConnection, error) {
	if connNotice == nil {
		lockNotice.Lock()
		defer lockNotice.Unlock()
		if connNotice == nil {
			Connect()
			if collection.Notice == nil {
				return nil, errors.New("failed to connect notice collection")
			}

			connNotice = &NoticeConnection{
				collection: collection.Notice,
			}
		}
	}

	return connNotice, nil
}

// Insert insert the notice document
func (conn *NoticeConnection) Insert(notice *types.Notice) error {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	res, err := conn.collection.InsertOne(ctx, &notice)
	if err != nil {
		return err
	}

	logger.Infof("insert notice collection successfully, %v", res.InsertedID)
	return nil
}

// Update insert the notice document
func (conn *NoticeConnection) Update(notice *types.Notice) error {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"_id": notice.ID}, bson.M{"$set": &notice}, op)
	if err != nil {
		return err
	}

	logger.Infof("insert notice collection successfully, %v", res.UpsertedID)
	return nil
}

// InsertMany insert the notice documents
func (conn *NoticeConnection) InsertMany(noticeList []*types.Notice) error {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	notices := []interface{}{}
	for _, notice := range noticeList {
		notices = append(notices, notice)
	}
	res, err := conn.collection.InsertMany(ctx, notices)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany notice collection successfully, %v", res.InsertedIDs)
	return nil
}

// Delete delete the notice ducument by id
func (conn *NoticeConnection) Delete(id string) error {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	logger.Infof("delete notice collection successfully, id:%s, DeletedCount:%d", id, res.DeletedCount)
	return nil
}

// DeleteMany delete the notice documents by query filter
func (conn *NoticeConnection) DeleteMany(query interface{}) error {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany notice collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

// Count query the number of all notice document in the mongodb
func (conn *NoticeConnection) Count() (int64, error) {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the notice collection successfully, %d", num)
	return num, nil
}

// QueryCount query the number of the notice document by query filter in the mongodb
func (conn *NoticeConnection) QueryCount(query interface{}) (int64, error) {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the notice collection successfully, %d", num)
	return num, nil
}

// Get get the notice by id
func (conn *NoticeConnection) Get(id string) (*types.Notice, error) {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	notice := &types.Notice{}
	if err := conn.collection.FindOne(ctx, bson.M{"_id": id}).Decode(notice); err != nil {
		return nil, err
	}
	return notice, nil
}

// GetByFilter get the notices by filter
func (conn *NoticeConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.Notice, error) {
	lockNotice.Lock()
	defer lockNotice.Unlock()

	notices := []*types.Notice{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"create_time": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		notice := &types.Notice{}
		if err = cursor.Decode(notice); err != nil {
			return nil, err
		}
		notices = append(notices, notice)
	}
	return notices, nil
}
