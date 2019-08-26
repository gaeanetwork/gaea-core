package mongodb

import (
	"errors"
	"sync"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//BlockInfoConnection for mongo db
type BlockInfoConnection struct {
	collection *mongo.Collection
}

var (
	connBlock *BlockInfoConnection
	lockBlock sync.Mutex
)

func GetBlockInfoConnection() (*BlockInfoConnection, error) {
	if connBlock == nil {
		lockBlock.Lock()
		defer lockBlock.Unlock()

		if connBlock == nil {
			Connect()
			if collection.BlockInfo == nil {
				return nil, errors.New("failed to connect block collection")
			}

			connBlock = &BlockInfoConnection{
				collection: collection.BlockInfo,
			}
		}
	}
	return connBlock, nil
}

func (conn *BlockInfoConnection) Insert(block *types.BlockInfo) error {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	res, err := conn.collection.InsertOne(ctx, &block)
	if err != nil {
		return err
	}

	logger.Infof("insert block collection successfully, %v", res.InsertedID)
	return nil
}

func (conn *BlockInfoConnection) Update(block *types.BlockInfo) error {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	op := options.Update()
	op.SetUpsert(true)

	res, err := conn.collection.UpdateOne(ctx, bson.M{"_id": block.ID}, &block, op)
	if err != nil {
		return err
	}

	logger.Infof("update block collection successfully, %v", res.UpsertedID)
	return nil
}

func (conn *BlockInfoConnection) InsertMany(blockList []*types.BlockInfo) error {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	blocks := []interface{}{}
	for _, block := range blockList {
		blocks = append(blocks, block)
	}
	res, err := conn.collection.InsertMany(ctx, blocks)
	if err != nil {
		return err
	}

	logger.Infof("InsertMany blocks collection successfully, %v", res.InsertedIDs)
	return nil
}

func (conn *BlockInfoConnection) Delete(blockhash string) error {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	res, err := conn.collection.DeleteOne(ctx, bson.D{{"blockhash", blockhash}})
	if err != nil {
		return err
	}

	logger.Infof("delete block collection successfully, id:%s, DeletedCount:%d", blockhash, res.DeletedCount)
	return nil
}

func (conn *BlockInfoConnection) DeleteMany(query interface{}) error {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	res, err := conn.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	logger.Infof("DeleteMany block collection successfully,  DeletedCount:%d", res.DeletedCount)
	return nil
}

func (conn *BlockInfoConnection) Count() (int64, error) {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	num, err := conn.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	logger.Infof("query all count of the block collection successfully, %d", num)
	return num, nil
}

func (conn *BlockInfoConnection) QueryCount(query interface{}) (int64, error) {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	num, err := conn.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	logger.Infof("query count of the block collection successfully, %d", num)
	return num, nil
}

func (conn *BlockInfoConnection) Get(blockhash string) (*types.BlockInfo, error) {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	block := &types.BlockInfo{}
	if err := conn.collection.FindOne(ctx, bson.D{{"blockhash", blockhash}}).Decode(block); err != nil {
		return nil, err
	}
	return block, nil
}

func (conn *BlockInfoConnection) GetByFilter(query interface{}, skip, limit int64) ([]*types.BlockInfo, error) {
	lockBlock.Lock()
	defer lockBlock.Unlock()

	blocks := []*types.BlockInfo{}

	cursor, err := conn.collection.Find(ctx, query, options.Find().SetSkip(skip).SetLimit(limit), options.Find().SetSort(bson.M{"timestamp": -1}))
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		block := &types.BlockInfo{}
		if err = cursor.Decode(block); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}
