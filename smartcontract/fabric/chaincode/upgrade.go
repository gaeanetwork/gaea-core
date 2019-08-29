package chaincode

import (
	"context"
	"fmt"

	protcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
)

const upgradeCmdName = "upgrade"

func upgrade(cfg *Config) error {
	cfg.CommandName = "upgrade"
	cf, err := InitCmdFactory(true, true, cfg)
	if err != nil {
		return err
	}
	defer cf.Close()

	env, err := chaincodeUpgrade(cf, cfg)
	if err != nil {
		return err
	}

	if env != nil {
		logger.Debug("Send signed envelope to orderer")
		err = cf.BroadcastClient.Send(env)
		return err
	}

	return nil
}

//upgrade the command via Endorser
func chaincodeUpgrade(cf *ChaincodeCmdFactory, cfg *Config) (*protcommon.Envelope, error) {
	spec, err := cfg.CreateChaincodeSpec()
	if err != nil {
		return nil, err
	}

	cds, err := getChaincodeDeploymentSpec(spec, false)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode code %s: %s", cfg.ChaincodeName, err)
	}

	creator, err := cf.Signer.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing identity for %s: %s", cf.Signer.GetIdentifier(), err)
	}

	prop, _, err := utils.CreateUpgradeProposalFromCDS(cfg.ChannelID, cds, creator, []byte{}, []byte{}, []byte{}, []byte{})
	if err != nil {
		return nil, fmt.Errorf("error creating proposal %s: %s", chainFuncName, err)
	}
	logger.Debugf("Get upgrade proposal for chaincode <%v>", spec.ChaincodeId)

	var signedProp *pb.SignedProposal
	signedProp, err = utils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return nil, fmt.Errorf("error creating signed proposal  %s: %s", chainFuncName, err)
	}

	// upgrade is currently only supported for one peer
	proposalResponse, err := cf.EndorserClients[0].ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, fmt.Errorf("error endorsing %s: %s", chainFuncName, err)
	}
	logger.Debugf("endorse upgrade proposal, get response <%v>", proposalResponse.Response)

	if proposalResponse != nil {
		// assemble a signed transaction (it's an Envelope message)
		env, err := utils.CreateSignedTx(prop, cf.Signer, proposalResponse)
		if err != nil {
			return nil, fmt.Errorf("could not assemble transaction, err %s", err)
		}
		logger.Debug("Get Signed envelope")
		return env, nil
	}

	return nil, nil
}
