package chain

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/scc/lscc"
	"github.com/pkg/errors"
)

const chaincodeLSCCName = "lscc"

// GetChaincodeDefinition get chaincode definition by name from lscc
func GetChaincodeDefinition(channelID, chaincodeName string) (ccprovider.ChaincodeDefinition, error) {
	chaincodeServer, err := getSystemChaincodeServer(chaincodeLSCCName, channelID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s chaincode server, channelID: %s", chaincodeLSCCName, channelID)
	}

	chaincodeDataStr, err := chaincodeServer.Query([]string{lscc.GETCCDATA, channelID, chaincodeName})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query system chaincode lscc, channelID: %s, chaincodeName: %s", channelID, chaincodeName)
	} else if chaincodeDataStr == "" {
		return nil, errors.Errorf("chaincode %s not found", chaincodeName)
	}

	chaincodeData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal([]byte(chaincodeDataStr), chaincodeData)
	if err != nil {
		return nil, errors.Wrapf(err, "chaincode %s has bad definition", chaincodeName)
	}

	return chaincodeData, nil
}
