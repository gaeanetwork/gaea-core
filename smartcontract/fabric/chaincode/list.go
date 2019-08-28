package chaincode

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	comm "gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
)

func list(cfg *Config, channelID string, getInstalledChaincodes, getInstantiatedChaincodes bool) ([]*types.ChaincodeInfo, error) {
	if !getInstalledChaincodes && getInstantiatedChaincodes && channelID == "" {
		return nil, errors.New("not specified channelID")
	}

	cfg.CommandName = "list"
	cfg.ChannelID = channelID
	cf, err := InitCmdFactory(true, false, cfg)
	if err != nil {
		return nil, err
	}
	defer cf.Close()

	creator, err := cf.Signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("Error serializing identity for %s: %s", cf.Signer.GetIdentifier(), err)
	}

	var prop *peer.Proposal
	if getInstalledChaincodes {
		prop, _, err = utils.CreateGetInstalledChaincodesProposal(creator)
	} else if getInstantiatedChaincodes {
		prop, _, err = utils.CreateGetChaincodesProposal(cfg.ChannelID, creator)
	} else {
		return nil, fmt.Errorf("Must explicitly specify \"--installed\" or \"--instantiated\"")
	}

	if err != nil {
		return nil, fmt.Errorf("Error creating proposal %s: %s", chainFuncName, err)
	}

	var signedProp *peer.SignedProposal
	signedProp, err = utils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return nil, fmt.Errorf("Error creating signed proposal  %s: %s", chainFuncName, err)
	}

	// list is currently only supported for one peer
	proposalResponse, err := cf.EndorserClients[0].ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, errors.Errorf("Error endorsing %s: %s", chainFuncName, err)
	}

	if proposalResponse.Response == nil {
		return nil, errors.Errorf("Proposal response had nil 'response'")
	}

	if proposalResponse.Response.Status != int32(common.Status_SUCCESS) {
		return nil, errors.Errorf("Bad response: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}

	cqr := &peer.ChaincodeQueryResponse{}
	err = proto.Unmarshal(proposalResponse.Response.Payload, cqr)
	if err != nil {
		return nil, err
	}

	infoList := []*types.ChaincodeInfo{}
	for _, cc := range cqr.Chaincodes {
		ccinfo := &types.ChaincodeInfo{
			ID:      comm.BytesToHex(cc.Id),
			Name:    cc.Name,
			Version: cc.Version,
			Path:    cc.Path,
			Input:   cc.Input,
			Escc:    cc.Escc,
			Vscc:    cc.Vscc,
		}
		infoList = append(infoList, ccinfo)
	}
	return infoList, nil
}
