package chaincode

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	pcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

const installCmdName = "install"

const installDesc = "Package the specified chaincode into a deployment spec and save it on the peer's path."

func install(cfg *Config, ccpackfile []byte) error {
	cfg.CommandName = "install"
	cf, err := InitCmdFactory(true, false, cfg)
	if err != nil {
		return err
	}
	defer cf.Close()

	ccpackmsg, _, err := GetPackageFromFile(ccpackfile)
	if err != nil {
		return err
	}

	return chaincodeInstall(ccpackmsg, cf)
}

func chaincodeInstall(msg proto.Message, cf *ChaincodeCmdFactory) error {
	creator, err := cf.Signer.Serialize()
	if err != nil {
		return fmt.Errorf("Error serializing identity for %s: %s", cf.Signer.GetIdentifier(), err)
	}

	prop, _, err := utils.CreateInstallProposalFromCDS(msg, creator)
	if err != nil {
		return fmt.Errorf("Error creating proposal  %s: %s", chainFuncName, err)
	}

	var signedProp *pb.SignedProposal
	signedProp, err = utils.GetSignedProposal(prop, cf.Signer)
	if err != nil {
		return fmt.Errorf("Error creating signed proposal  %s: %s", chainFuncName, err)
	}

	// install is currently only supported for one peer
	proposalResponse, err := cf.EndorserClients[0].ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return fmt.Errorf("Error endorsing %s: %s", chainFuncName, err)
	}

	if proposalResponse != nil {
		if proposalResponse.Response.Status != int32(pcommon.Status_SUCCESS) {
			return errors.Errorf("Bad response: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
		}
		logger.Infof("Installed remotely %v", proposalResponse)
	} else {
		return errors.New("Error during install: received nil proposal response")
	}

	return nil
}

//GetPackageFromFile get the chaincode package from file and the extracted ChaincodeDeploymentSpec
func GetPackageFromFile(packfile []byte) (proto.Message, *pb.ChaincodeDeploymentSpec, error) {
	ccpack, err := ccprovider.GetCCPackage(packfile)
	if err != nil {
		return nil, nil, err
	}

	o := ccpack.GetPackageObject()

	cds, ok := o.(*pb.ChaincodeDeploymentSpec)
	if !ok || cds == nil {
		env, ok := o.(*pcommon.Envelope)
		if !ok || env == nil {
			return nil, nil, fmt.Errorf("error extracting valid chaincode package")
		}

		_, sCDS, err := ccpackage.ExtractSignedCCDepSpec(env)
		if err != nil {
			return nil, nil, fmt.Errorf("error extracting valid signed chaincode package(%s)", err)
		}

		cds, err = utils.GetChaincodeDeploymentSpec(sCDS.ChaincodeDeploymentSpec, platformRegistry)
		if err != nil {
			return nil, nil, fmt.Errorf("error extracting chaincode deployment spec(%s)", err)
		}
	}

	return o, cds, nil
}
