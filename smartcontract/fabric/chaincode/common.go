package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/pkg/errors"
)

func chaincodeInvokeOrQuery(invoke bool, cf *ChaincodeCmdFactory, cfg *Config) (result string, err error) {
	if err := checkChaincodeCmdParams(cfg); err != nil {
		return "", err
	}

	spec, err := cfg.CreateChaincodeSpec()
	if err != nil {
		return "", err
	}
	// call with empty txid to ensure production code generates a txid.
	// otherwise, tests can explicitly set their own txid
	txID := ""

	proposalResp, err := ChaincodeInvokeOrQuery(
		spec,
		cfg.ChannelID,
		txID,
		invoke,
		cf.Signer,
		cf.Certificate,
		cf.EndorserClients,
		cf.DeliverClients,
		cf.BroadcastClient,
		cfg)

	if err != nil {
		return "", errors.Errorf("%s - proposal response: %v", err, proposalResp)
	}

	if invoke {
		pRespPayload, err := putils.GetProposalResponsePayload(proposalResp.Payload)
		if err != nil {
			return "", errors.WithMessage(err, "error while unmarshaling proposal response payload")
		}
		ca, err := putils.GetChaincodeAction(pRespPayload.Extension)
		if err != nil {
			return "", errors.WithMessage(err, "error while unmarshaling chaincode action")
		}
		if proposalResp.Endorsement == nil {
			return "", errors.Errorf("endorsement failure during invoke. response: %v", proposalResp.Response)
		}

		if ca.Response.Status != shim.OK {
			return "", errors.Errorf("Response status is not 200. response: %v", proposalResp.Response)
		}
		return string(ca.Response.Payload), nil
	}

	if proposalResp == nil {
		return "", errors.New("error during query: received nil proposal response")
	}
	if proposalResp.Endorsement == nil {
		return "", errors.Errorf("endorsement failure during query. response: %v", proposalResp.Response)
	}

	if cfg.ChaincodeQueryRaw && cfg.ChaincodeQueryHex {
		return "", fmt.Errorf("options --raw (-r) and --hex (-x) are not compatible")
	}
	if cfg.ChaincodeQueryRaw {
		return fmt.Sprint(proposalResp.Response.Payload), nil
	}
	if cfg.ChaincodeQueryHex {
		return fmt.Sprintf("%x", proposalResp.Response.Payload), nil
	}

	return string(proposalResp.Response.Payload), nil
}

func checkChaincodeCmdParams(cfg *Config) error {
	// we need chaincode name for everything, including deploy
	if cfg.ChaincodeName == common.UndefinedParamValue {
		return errors.Errorf("must supply value for %s name parameter", chainFuncName)
	}

	if cfg.CommandName == instantiateCmdName || cfg.CommandName == installCmdName ||
		cfg.CommandName == upgradeCmdName || cfg.CommandName == packageCmdName {
		if cfg.ChaincodeVersion == common.UndefinedParamValue {
			return errors.Errorf("chaincode version is not provided for %s", cfg.CommandName)
		}

		escc, vscc := cfg.ESCC, cfg.VSCC
		if escc != common.UndefinedParamValue {
			logger.Infof("Using escc %s", escc)
		} else {
			logger.Info("Using default escc")
			escc = "escc"
		}

		if vscc != common.UndefinedParamValue {
			logger.Infof("Using vscc %s", vscc)
		} else {
			logger.Info("Using default vscc")
			vscc = "vscc"
		}

		policy := cfg.Policy
		if policy != common.UndefinedParamValue {
			p, err := cauthdsl.FromString(policy)
			if err != nil {
				return errors.Errorf("invalid policy %s", policy)
			}
			cfg.PolicyMarshalled = putils.MarshalOrPanic(p)
		}

		collectionsConfigFile := cfg.CollectionsConfigFile
		if collectionsConfigFile != common.UndefinedParamValue {
			var err error
			cfg.CollectionConfigBytes, err = getCollectionConfigFromFile(collectionsConfigFile)
			if err != nil {
				return errors.WithMessage(err, fmt.Sprintf("invalid collection configuration in file %s", collectionsConfigFile))
			}
		}
	}

	// TODO check cfg.CHaincodeInput
	return nil
}
