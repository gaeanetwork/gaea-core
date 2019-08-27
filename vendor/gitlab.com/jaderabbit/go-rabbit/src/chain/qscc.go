package chain

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	sdkcommon "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"gitlab.com/jaderabbit/go-rabbit/api/models"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/sdk"
	comm "gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
)

const chaincodeQSCCName = "qscc"

// GetTransactionByID get transaction by id
func GetTransactionByID(txID, channelID string) (*types.Transaction, error) {
	sdkLedger, err := sdk.GetLedger(channelID)
	if err != nil {
		return nil, err
	}

	processedTx, err := sdkLedger.QueryTransaction(txID)
	if err != nil {
		return nil, err
	}

	payload, err := utils.UnmarshalPayload(processedTx.TransactionEnvelope.Payload)
	if err != nil {
		return nil, err
	}

	tx, err := utils.GetTransaction(payload.Data)
	if err != nil {
		return nil, err
	}
	channelHeader := &common.ChannelHeader{}
	err = proto.Unmarshal(payload.Header.ChannelHeader, channelHeader)
	if err != nil {
		return nil, err
	}

	transaction := types.Transaction{}
	transaction.TimeStamp = time.Unix(channelHeader.Timestamp.Seconds, 0)
	transaction.TxID = channelHeader.TxId

	txAction := tx.Actions[0]
	actionPayload, action, err := utils.GetPayloads(txAction)
	if err != nil {
		return nil, err
	}

	chaincodeProposalPayload, err := utils.GetChaincodeProposalPayload(actionPayload.ChaincodeProposalPayload)
	if err != nil {
		return nil, err
	}

	cis := &peer.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(chaincodeProposalPayload.Input, cis)
	if err != nil {
		return nil, err
	}

	inputByte, err := proto.Marshal(cis.ChaincodeSpec.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ChaincodeSpec.ChaincodeSpec, err:%s", err.Error())
	}

	transaction.ChannelID = cis.ChaincodeSpec.ChaincodeId
	transaction.Input = string(inputByte)
	transaction.Response = action.Response.String()

	txRwSet := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(action.Results, txRwSet)
	if err != nil {
		return nil, err
	}

	for _, nsrs := range txRwSet.NsRwset {
		kvset := &kvrwset.KVRWSet{}
		err = proto.Unmarshal(nsrs.Rwset, kvset)
		if err != nil {
			return nil, err
		}

		for _, write := range kvset.Writes {
			transaction.KVWSet = append(transaction.KVWSet, &types.KVWrite{
				Key:      write.Key,
				Value:    string(write.Value),
				IsDelete: write.IsDelete,
			})
		}
	}

	return &transaction, nil
}

// GetChainInfo get chain info from the fabric
func GetChainInfo(channelID string) (*sdkcommon.BlockchainInfo, error) {
	sdkLedger, err := sdk.GetLedger(channelID)
	if err != nil {
		return nil, err
	}

	response, err := sdkLedger.QueryInfo()
	if err != nil {
		return nil, err
	}

	return response.BCI, nil
}

// GetBlockByHash get block by hash and channelID
func GetBlockByHash(channelID, blockHash string) (*common.Block, error) {
	sdkLedger, err := sdk.GetLedger(channelID)
	if err != nil {
		return nil, err
	}

	block, err := sdkLedger.QueryBlockByHash([]byte(blockHash))
	if err != nil {
		return nil, err
	}

	return ConvertSDKBlock(block)
}

// GetBlockByHeight get block by height and channelID
func GetBlockByHeight(channelID string, height uint64) (*common.Block, error) {
	sdkLedger, err := sdk.GetLedger(channelID)
	if err != nil {
		return nil, err
	}

	block, err := sdkLedger.QueryBlock(height)
	if err != nil {
		return nil, err
	}
	return ConvertSDKBlock(block)
}

// GetTxsInBlock get transaction information from the block
func GetTxsInBlock(block *common.Block) ([]*types.TxInfo, error) {
	buf := proto.NewBuffer(nil)
	blockData := block.Data
	txs := make([]*types.TxInfo, 0)

	if err := buf.EncodeVarint(uint64(len(blockData.Data))); err != nil {
		return nil, err
	}

	for _, txEnvelopeBytes := range blockData.Data {
		txEnvelope, err := utils.GetEnvelopeFromBlock(txEnvelopeBytes)
		if err != nil {
			return nil, err
		}

		txPayload, err := utils.GetPayload(txEnvelope)
		if err != nil {
			return nil, err
		}

		sh, err := utils.GetSignatureHeader(txPayload.Header.SignatureHeader)
		if err != nil {
			return nil, err
		}

		mspID, err := GetID(sh.Creator)
		if err != nil {
			return nil, fmt.Errorf("failed to get msp id of transaction, err:%s", err.Error())
		}

		userID := "system"
		user, err := models.QueryUserByMspID(mspID)
		if err != nil {
			logger.Warnf("failed to get user by mspid, mspID:%s, err:%s", mspID, err.Error)
		} else {
			userID = user.ID
		}

		chdr, err := utils.UnmarshalChannelHeader(txPayload.Header.ChannelHeader)
		if err != nil {
			return nil, err
		}

		chaincodeHeaderExtension, err := utils.GetChaincodeHeaderExtension(txPayload.Header)
		if err != nil {
			return nil, err
		}

		chaincodeName := ""
		if chaincodeHeaderExtension.ChaincodeId != nil {
			chaincodeName = chaincodeHeaderExtension.ChaincodeId.Name
		}

		id := fmt.Sprintf("%s-%s", chdr.ChannelId, chdr.TxId)

		tx := &types.TxInfo{
			ID:            id,
			TxID:          chdr.TxId,
			BlockHash:     comm.BytesToHex(block.Header.Hash()),
			Number:        block.Header.Number,
			ChannelID:     chdr.ChannelId,
			ChaincodeName: chaincodeName,
			UserID:        userID,
			Timestamp:     chdr.Timestamp.Seconds,
		}

		txs = append(txs, tx)
	}
	return txs, nil
}

// ConvertSDKBlock get sdk transaction from the ledger, marshal to byte array, then unmarshal to common block
func ConvertSDKBlock(sdkBlock *sdkcommon.Block) (*common.Block, error) {
	byteBlock, err := utils.Marshal(sdkBlock)
	if err != nil {
		return nil, err
	}

	return utils.UnmarshalBlock(byteBlock)
}
