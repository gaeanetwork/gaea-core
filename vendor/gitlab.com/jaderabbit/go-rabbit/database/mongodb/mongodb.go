package mongodb

import (
	"context"
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionMap struct {
	TxInfo                *mongo.Collection
	BlockInfo             *mongo.Collection
	Statistics            *mongo.Collection
	Notice                *mongo.Collection
	ChainTPS              *mongo.Collection
	Asset                 *mongo.Collection
	ActionLog             *mongo.Collection
	ChaincodeInvokeResult *mongo.Collection
}

var (
	// db the database in the mongodb
	db *mongo.Database
	// collection the collections int the database
	collection collectionMap
	ctx        context.Context

	logger = flogging.MustGetLogger("database.mongodb")
	once   sync.Once
)

// Connect connect the mongodb, and instantiate DB and Collection
func Connect() {
	once.Do(func() {
		ctx := context.Background()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			logger.Warningf("mongo connect error:%s", err.Error())
			return
		}

		db = client.Database("go-rabbit")
		collection.TxInfo = db.Collection("TxInfo")
		collection.BlockInfo = db.Collection("BlockInfo")
		collection.Statistics = db.Collection("Statistics")
		collection.Notice = db.Collection("Notice")
		collection.ChainTPS = db.Collection("ChainTPS")
		collection.Asset = db.Collection("Asset")
		collection.ActionLog = db.Collection("ActionLog")
		collection.ChaincodeInvokeResult = db.Collection("ChaincodeInvokeResult")
	})
}
