package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ChaincodeInvokeResultConnection for mongo db
type ChaincodeInvokeResultConnection struct {
	collection *mongo.Collection
}

var (
	connChaincodeInvokeResult *ChaincodeInvokeResultConnection
	lockChaincodeInvokeResult sync.Mutex
)

// GetChaincodeInvokeResultConnection get the ChaincodeInvokeResult connection from the mongodb
func GetChaincodeInvokeResultConnection() (*ChaincodeInvokeResultConnection, error) {
	if connChaincodeInvokeResult == nil {
		lockChaincodeInvokeResult.Lock()
		defer lockChaincodeInvokeResult.Unlock()
		if connChaincodeInvokeResult == nil {
			Connect()
			if collection.ChaincodeInvokeResult == nil {
				return nil, errors.New("failed to connect ChaincodeInvokeResult collection")
			}

			connChaincodeInvokeResult = &ChaincodeInvokeResultConnection{
				collection: collection.ChaincodeInvokeResult,
			}
		}
	}

	return connChaincodeInvokeResult, nil
}

// Insert insert the ChaincodeInvokeResult document
func (conn *ChaincodeInvokeResultConnection) Insert(ccIR *types.ChaincodeInvokeResult) error {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	res, err := conn.collection.InsertOne(ctx, &ccIR)
	if err != nil {
		return err
	}

	logger.Debugf("insert ChaincodeInvokeResult collection successfully, %v", res.InsertedID)
	return nil
}

// InsertMany insert the ChaincodeInvokeResult documents
func (conn *ChaincodeInvokeResultConnection) InsertMany(ccIRList []*types.ChaincodeInvokeResult) error {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	ccIRs := []interface{}{}
	for _, ccIR := range ccIRList {
		ccIRs = append(ccIRs, ccIR)
	}
	res, err := conn.collection.InsertMany(ctx, ccIRs)
	if err != nil {
		return err
	}

	logger.Debugf("InsertMany ChaincodeInvokeResult collection successfully, %v", res.InsertedIDs)
	return nil
}

// Delete delete the ChaincodeInvokeResult ducument by id
func (conn *ChaincodeInvokeResultConnection) Delete(id string) error {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	logger.Infof("delete ChaincodeInvokeResult collection successfully, id:%s, DeletedCount:%d", id, res.DeletedCount)
	return nil
}

// DeleteMany delete the ChaincodeInvokeResult documents by query filter
func (conn *ChaincodeInvokeResultConnection) DeleteMany(query interface{}) error {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany ChaincodeInvokeResult collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

// Count query the number of all ChaincodeInvokeResult document in the mongodb
func (conn *ChaincodeInvokeResultConnection) Count() (int64, error) {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the ChaincodeInvokeResult collection successfully, %d", num)
	return num, nil
}

// QueryCount query the number of the ChaincodeInvokeResult document by query filter in the mongodb
func (conn *ChaincodeInvokeResultConnection) QueryCount(query interface{}) (int64, error) {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Debugf("query count of the ChaincodeInvokeResult collection successfully, %d", num)
	return num, nil
}

// Get get the ChaincodeInvokeResult by id
func (conn *ChaincodeInvokeResultConnection) Get(id string) (*types.ChaincodeInvokeResult, error) {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	ccIR := &types.ChaincodeInvokeResult{}
	if err := conn.collection.FindOne(ctx, bson.M{"_id": id}).Decode(ccIR); err != nil {
		return nil, err
	}
	return ccIR, nil
}

// GetByFilter get the ChaincodeInvokeResult by filter
func (conn *ChaincodeInvokeResultConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.ChaincodeInvokeResult, error) {
	lockChaincodeInvokeResult.Lock()
	defer lockChaincodeInvokeResult.Unlock()

	ccIRs := []*types.ChaincodeInvokeResult{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"timestamp": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		ccIR := &types.ChaincodeInvokeResult{}
		if err = cursor.Decode(ccIR); err != nil {
			return nil, err
		}
		ccIRs = append(ccIRs, ccIR)
	}
	return ccIRs, nil
}
