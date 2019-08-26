package chain

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"gitlab.com/jaderabbit/go-rabbit/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTransactionList get transaction list from mongodb by channelID
func GetTransactionList(channelID, chaincodeName string, pageIndex, pageSize uint32) ([]*types.TxInfo, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return nil, err
	}

	skip := (pageIndex - 1) * pageSize

	query := make(bson.M)
	query["channel"] = channelID

	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	return conn.GetByFilter(query, int64(skip), int64(pageSize))
}

// GetTransactionListByFilter get transaction list from mongodb by channelID
func GetTransactionListByFilter(query interface{}, pageIndex, pageSize uint32) ([]*types.TxInfo, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return nil, err
	}

	skip := (pageIndex - 1) * pageSize
	return conn.GetByFilter(query, int64(skip), int64(pageSize))
}

// GetTransactionNum get the number of all channel transaction from the mongodb
func GetTransactionNum() (int64, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return 0, err
	}

	return conn.Count()
}

// GetTransactionNumByChannelID get the number of the transaction from the mongodb by channelID
func GetTransactionNumByChannelID(channelID, chaincodeName string) (int64, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return 0, err
	}
	query := bson.M{"channel": channelID}
	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}
	return conn.QueryCount(query)
}

// GetTransactionNumByFilter get the number of the transaction from the mongodb by query
func GetTransactionNumByFilter(query interface{}) (int64, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return 0, err
	}
	return conn.QueryCount(query)
}

// GetTodayTransactionNum get the number of today transaction from the mongodb by channelID
func GetTodayTransactionNum(channelID, chaincodeName string) (int64, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return 0, err
	}

	nowTime := time.Now()
	todayTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
	filter := bson.M{"$gte": todayTime.Unix()}
	query := bson.M{"timestamp": filter, "channel": channelID}

	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	return conn.QueryCount(query)
}

// GetLastHourTransactionNum get the number of last hour transaction from the mongodb by channelID
func GetLastHourTransactionNum(channelID, chaincodeName string) (int64, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return 0, err
	}

	nowTime := time.Now().Add(-1 * time.Hour)
	filter := bson.M{"$gte": nowTime.Unix()}
	query := bson.M{"timestamp": filter, "channel": channelID}

	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	return conn.QueryCount(query)
}

// GetBlockNum get the number of all channel block from the mongodb
func GetBlockNum() (int64, error) {
	conn, err := mongodb.GetBlockInfoConnection()
	if err != nil {
		return 0, err
	}

	return conn.Count()
}

// GetLastHourBlockNum get the number of last hour block from mongodb by channelID
func GetLastHourBlockNum(channelID, chaincodeName string) (int64, error) {
	conn, err := mongodb.GetBlockInfoConnection()
	if err != nil {
		return 0, err
	}

	nowTime := time.Now().Add(-1 * time.Hour)
	filter := bson.M{"$gte": nowTime.Unix()}
	query := bson.M{"timestamp": filter, "channel": channelID}

	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	return conn.QueryCount(query)
}

// GetBlockNumByChannelID get the number of block from mongodb by channelID
func GetBlockNumByChannelID(channelID, chaincodeName string) (int64, error) {
	conn, err := mongodb.GetBlockInfoConnection()
	if err != nil {
		return 0, err
	}

	query := bson.M{"channel": channelID}

	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	return conn.QueryCount(query)
}

// GetStatistics get the statistics data by day in a month
func GetStatistics(channelID, chaincodeName string, year, month int) (map[string]*types.Statistics, error) {
	conn, err := mongodb.GetStatisticsConnection()
	if err != nil {
		return nil, err
	}

	if len(channelID) == 0 {
		return nil, errors.New("channelID is empty")
	}

	query := bson.M{"channel": channelID, "year": year, "month": month}
	if len(chaincodeName) > 0 {
		query["chaincode"] = chaincodeName
	}

	list, err := conn.GetByFilter(query, 0, 32*24)
	if err != nil {
		return nil, err
	}

	statisticsTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	mapStatistics := make(map[string]*types.Statistics)

	for statisticsTime.Year() == int(year) && statisticsTime.Month() == time.Month(month) {
		mapStatistics[fmt.Sprint(statisticsTime.Day())] = &types.Statistics{
			Year:        year,
			Month:       month,
			Day:         statisticsTime.Day(),
			ChannelID:   channelID,
			TxNumber:    0,
			BlockNumber: 0,
		}
		statisticsTime = statisticsTime.Add(24 * time.Hour)
	}

	for _, statistics := range list {
		key := fmt.Sprint(statistics.Day)
		s := mapStatistics[key]

		s.TxNumber += statistics.TxNumber
		s.BlockNumber += statistics.BlockNumber
	}

	return mapStatistics, nil
}

func GetTPSByChannelID(channelID string) (int, error) {
	chainTPS, err := mongodb.GetChainTPSInfoConnection()
	if err != nil {
		return 0, err
	}

	query := bson.M{}
	if len(channelID) > 0 {
		query = bson.M{"channelid": channelID}
	}

	tpss, err := chainTPS.GetByFilter(query, 0, 1)
	if err != nil {
		return 0, err
	}

	if len(tpss) == 0 {
		return 0, nil
	}

	return tpss[0].TPS, nil
}

func StatisticsTxGroupByUser(channelID string, beginTime, endTime int64) ([]*types.UserGroupTx, error) {
	conn, err := mongodb.GetTxInfoConnection()
	if err != nil {
		return nil, err
	}

	filterTime := bson.M{"$gte": beginTime, "$lte": endTime}
	filterChannel := bson.M{"timestamp": filterTime, "channel": channelID, "user_id": bson.M{"$ne": "system"}}
	filter := []bson.M{{"$match": filterChannel},
		{"$group": bson.M{"_id": "$user_id", "count": bson.M{"$sum": 1}}},
	}

	return conn.Aggregate(filter)
}
