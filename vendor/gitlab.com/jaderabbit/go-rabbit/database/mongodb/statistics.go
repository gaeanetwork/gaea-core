package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//StatisticsConnection for mongo db
type StatisticsConnection struct {
	collection *mongo.Collection
}

var (
	connStatistic *StatisticsConnection
	lockStatistic sync.Mutex
)

func GetStatisticsConnection() (*StatisticsConnection, error) {
	if connStatistic == nil {
		lockStatistic.Lock()
		defer lockStatistic.Unlock()
		if connStatistic == nil {
			Connect()
			if collection.Statistics == nil {
				return nil, errors.New("failed to connect Statistics collection")
			}

			connStatistic = &StatisticsConnection{
				collection: collection.Statistics,
			}
		}
	}

	return connStatistic, nil
}

func (conn *StatisticsConnection) Insert(statistics *types.Statistics) error {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	res, err := conn.collection.InsertOne(ctx, &statistics)
	if err != nil {
		return err
	}

	logger.Infof("insert statistics collection successfully, %v", res.InsertedID)
	return nil
}

func (conn *StatisticsConnection) Update(statistics *types.Statistics) error {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"_id": statistics.ID}, bson.M{"$set": &statistics}, op)
	if err != nil {
		return err
	}

	logger.Infof("update statistics collection successfully, %v", res.UpsertedID)
	return nil
}

func (conn *StatisticsConnection) InsertMany(statisticsList []*types.Statistics) error {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	statisticsArray := []interface{}{}
	for _, statistics := range statisticsList {
		statisticsArray = append(statisticsArray, statistics)
	}

	res, err := conn.collection.InsertMany(ctx, statisticsArray)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany statistics collection successfully, %v", res.InsertedIDs)
	return nil
}

func (conn *StatisticsConnection) Delete(id string) error {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	logger.Infof("Delete statistics collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *StatisticsConnection) DeleteMany(query interface{}) error {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany statistics collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *StatisticsConnection) Count() (int64, error) {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the statistics collection successfully, %d", num)
	return num, nil
}

func (conn *StatisticsConnection) QueryCount(query interface{}) (int64, error) {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the statistics collection successfully, %d", num)
	return num, nil
}

func (conn *StatisticsConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.Statistics, error) {
	lockStatistic.Lock()
	defer lockStatistic.Unlock()

	statisticsArray := []*types.Statistics{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"timestamp": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		statistics := &types.Statistics{}
		if err = cursor.Decode(statistics); err != nil {
			return nil, err
		}
		statisticsArray = append(statisticsArray, statistics)
	}
	return statisticsArray, nil
}
