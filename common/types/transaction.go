package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/peer"
)

// Transaction the transaction infomation from the fabric ledger
type Transaction struct {
	TxID      string            `json:"txID"`
	TimeStamp time.Time         `json:"timeStamp"`
	ChannelID *peer.ChaincodeID `json:"channelID"`
	Input     string            `json:"input"`
	Response  string            `json:"response"`
	KVWSet    []*KVWrite        `json:"kvwset"`
}

// KVWrite captures a write (update/delete) operation performed during transaction simulation
type KVWrite struct {
	Key      string `json:"key,omitempty"`
	IsDelete bool   `json:"is_delete,omitempty"`
	Value    string `json:"value,omitempty"`
}

// ChaincodeInvokeResult chaincode invoke result
type ChaincodeInvokeResult struct {
	ID               string              `bson:"_id" json:"id"`
	ChannelID        string              `bson:"channel_id" json:"channel_id"`
	ChaincodeName    string              `bson:"chaincode_name" json:"chaincode_name"`
	InputArgs        []string            `bson:"input_args" json:"input_args"`
	TransactionID    fab.TransactionID   `bson:"tx_id" json:"tx_id"`
	TxValidationCode pb.TxValidationCode `bson:"tx_validation_code" json:"tx_validation_code"`
	ChaincodeStatus  int32               `bson:"chaincode_status" json:"chaincode_status"`
	Payload          []byte              `bson:"payload" json:"payload"`
	ErrorInfo        string              `bson:"error_info" json:"error_info"`
	Timestamp        int64               `bson:"timestamp" json:"timestamp,omitempty"`
}

func NewChaincodeInvokeResult(args []string) *ChaincodeInvokeResult {
	return &ChaincodeInvokeResult{
		ID:        uuid.New().String(),
		InputArgs: args,
		Timestamp: time.Now().Unix(),
	}
}
