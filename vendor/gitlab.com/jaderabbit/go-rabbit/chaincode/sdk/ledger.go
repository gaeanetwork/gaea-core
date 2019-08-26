package sdk

import (
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

var (
	mapLedger   map[string]*Ledger
	ledgerMutex sync.Mutex
)

type Ledger struct {
	ChannelID string
	client    *ledger.Client
}

func GetLedger(channelID string) (*Ledger, error) {
	ledgerMutex.Lock()
	defer ledgerMutex.Unlock()
	if mapLedger == nil {
		mapLedger = make(map[string]*Ledger)
	}

	sdkLedger, ok := mapLedger[channelID]
	if ok {
		return sdkLedger, nil
	}

	sdkLedger, err := NewLedger(channelID)
	if err != nil {
		return nil, err
	}

	mapLedger[channelID] = sdkLedger
	return sdkLedger, nil
}

func NewLedger(channelID string) (*Ledger, error) {
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to new fabsdk, err:%s", err.Error())
	}

	//prepare required contexts
	channelClientCtx := sdk.ChannelContext(channelID, fabsdk.WithUser(DefaultMspUserName), fabsdk.WithOrg(orgName))

	ledgerClient, err := ledger.New(channelClientCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to new ledger, err:%s", err.Error())
	}

	return &Ledger{
		ChannelID: channelID,
		client:    ledgerClient,
	}, nil
}

// QueryTransaction queries the ledger for processed transaction by transaction ID.
//  Parameters:
//  txID is required transaction ID
//  options hold optional request options
//
//  Returns:
//  processed transaction information
func (l *Ledger) QueryTransaction(txID string, options ...ledger.RequestOption) (*pb.ProcessedTransaction, error) {
	return l.client.QueryTransaction(fab.TransactionID(txID), options...)
}

// QueryInfo queries for various useful blockchain information on this channel such as block height and current block hash.
//  Parameters:
//  options are optional request options
//
//  Returns:
//  blockchain information
func (l *Ledger) QueryInfo(options ...ledger.RequestOption) (*fab.BlockchainInfoResponse, error) {
	return l.client.QueryInfo(options...)
}

// QueryBlock queries the ledger for Block by block number.
//  Parameters:
//  blockNumber is required block number(ID)
//  options hold optional request options
//
//  Returns:
//  block information
func (l *Ledger) QueryBlock(blockNumber uint64, options ...ledger.RequestOption) (*common.Block, error) {
	return l.client.QueryBlock(blockNumber, options...)
}

// QueryBlockByHash queries the ledger for block by block hash.
//  Parameters:
//  blockHash is required block hash
//  options hold optional request options
//
//  Returns:
//  block information
func (l *Ledger) QueryBlockByHash(blockHash []byte, options ...ledger.RequestOption) (*common.Block, error) {
	return l.client.QueryBlockByHash(blockHash, options...)
}

// QueryBlockByTxID queries for block which contains a transaction.
//  Parameters:
//  txID is required transaction ID
//  options hold optional request options
//
//  Returns:
//  block information
func (l *Ledger) QueryBlockByTxID(txID string, options ...ledger.RequestOption) (*common.Block, error) {
	return l.client.QueryBlockByTxID(fab.TransactionID(txID), options...)
}

// QueryConfig queries for channel configuration.
//  Parameters:
//  options hold optional request options
//
//  Returns:
//  channel configuration information
func (l *Ledger) QueryConfig(options ...ledger.RequestOption) (fab.ChannelCfg, error) {
	return l.client.QueryConfig(options...)
}
