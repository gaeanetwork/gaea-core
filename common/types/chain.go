package types

import (
	"errors"
	"fmt"
)

//TxInfo for statistics
type TxInfo struct {
	ID            string `bson:"_id" json:"id,omitempty"`
	TxID          string `bson:"txid" json:"txid,omitempty"`
	BlockHash     string `bson:"blockhash" json:"blockhash,omitempty"`
	Number        uint64 `bson:"number" json:"number,omitempty"`
	ChannelID     string `bson:"channel" json:"channel,omitempty"`
	ChaincodeName string `bson:"chaincode" json:"chaincode,omitempty"`
	UserID        string `bson:"user_id" json:"user_id,omitempty"`
	Timestamp     int64  `bson:"timestamp" json:"timestamp,omitempty"`
}

type UserGroupTx struct {
	UserID string `bson:"_id" json:"user_id,omitempty"`
	Count  int    `bson:"count" json:"count,omitempty"`
}

type BlockInfo struct {
	ID            string `bson:"_id" json:"id,omitempty"`
	BlockHash     string `bson:"blockhash" json:"blockhash,omitempty"`
	Number        uint64 `bson:"number" json:"number,omitempty"`
	PreviousHash  string `bson:"previoushash" json:"prehash,omitempty"`
	ChannelID     string `bson:"channel" json:"channel,omitempty"`
	ChaincodeName string `bson:"chaincode" json:"chaincode,omitempty"`
	Timestamp     int64  `bson:"timestamp" json:"timestamp,omitempty"`
}

type Statistics struct {
	ID            string `bson:"_id" json:"id,omitempty"`
	Year          int    `bson:"year" json:"year,omitempty"`
	Month         int    `bson:"month" json:"month,omitempty"`
	Day           int    `bson:"day" json:"day,omitempty"`
	Hour          int    `bson:"hour" json:"hour,omitempty"`
	ChannelID     string `bson:"channel" json:"channel,omitempty"`
	TxNumber      int    `bson:"txnum" json:"txnum,omitempty"`
	BlockNumber   int    `bson:"blocknum" json:"blocknum,omitempty"`
	ChaincodeName string `bson:"chaincode" json:"chaincode,omitempty"`
}

type ChainTPS struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	ChannelID string `bson:"channelid" json:"channelid,omitempty"`
	TPS       int    `bson:"tps" json:"tps,omitempty"`
}

// ChaincodeInfo contains general information about an installed/instantiated chaincode
type ChaincodeInfo struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
	Input   string `json:"input,omitempty"`
	Escc    string `json:"escc,omitempty"`
	Vscc    string `json:"vscc,omitempty"`
}

// NewStatistics new a statistic by channelID, chaincodeName, year, month, day, hour
func NewStatistics(channelID, chaincodeName string, year, month, day, hour int) (*Statistics, error) {
	if len(channelID) == 0 {
		return nil, errors.New("not specified channelID")
	}

	if len(chaincodeName) == 0 {
		return nil, errors.New("not specified chaincodeName")
	}

	if year <= 0 {
		return nil, fmt.Errorf("year must be a positive integer, year:%d", year)
	}

	if month <= 0 {
		return nil, fmt.Errorf("month must be a positive integer, month:%d", month)
	}

	if day <= 0 {
		return nil, fmt.Errorf("day must be a positive integer, day:%d", day)
	}

	id := fmt.Sprintf("%s-%s-%d-%d-%d-%d", channelID, chaincodeName, year, month, day, hour)
	return &Statistics{
		ID:            id,
		Year:          year,
		Month:         month,
		Day:           day,
		Hour:          hour,
		ChannelID:     channelID,
		ChaincodeName: chaincodeName,
	}, nil
}
