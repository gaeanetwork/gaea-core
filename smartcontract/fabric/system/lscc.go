package system

import (
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/scc/lscc"
	"github.com/pkg/errors"
)

const chaincodeLSCCName = "lscc"

// GetChaincodeDefinition get chaincode definition by name from lscc
func GetChaincodeDefinition(channelID, chaincodeName string) (ccprovider.ChaincodeDefinition, error) {
	chaincodeService := &fabric.Chaincode{}

	chaincodeDataBytes, err := chaincodeService.Query(chaincodeLSCCName, []string{lscc.GETCCDATA, channelID, chaincodeName})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query system chaincode lscc, channelID: %s, chaincodeName: %s", channelID, chaincodeName)
	} else if len(chaincodeDataBytes) == 0 {
		return nil, errors.Errorf("chaincode %s not found", chaincodeName)
	}

	chaincodeData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(chaincodeDataBytes, chaincodeData)
	if err != nil {
		return nil, errors.Wrapf(err, "chaincode %s has bad definition", chaincodeName)
	}

	return chaincodeData, nil
}
