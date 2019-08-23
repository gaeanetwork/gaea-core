package chaincode

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/cmd"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/client"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

func chaincodeInvokeOrQuery(invoke bool, cf *cmd.ChaincodeCmdFactory, cfg *Config) (result string, err error) {
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

	proposalResp, err := InvokeOrQuery(
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
		pRespPayload, err := utils.GetProposalResponsePayload(proposalResp.Payload)
		if err != nil {
			return "", errors.WithMessage(err, "error while unmarshaling proposal response payload")
		}
		ca, err := utils.GetChaincodeAction(pRespPayload.Extension)
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
			cfg.PolicyMarshalled = utils.MarshalOrPanic(p)
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

// InvokeOrQuery invokes or queries the chaincode. If successful, the
// INVOKE form prints the ProposalResponse to STDOUT, and the QUERY form prints
// the query result on STDOUT. A command-line flag (-r, --raw) determines
// whether the query result is output as raw bytes, or as a printable string.
// The printable form is optionally (-x, --hex) a hexadecimal representation
// of the query response. If the query response is NIL, nothing is output.
//
// NOTE - Query will likely go away as all interactions with the endorser are
// Proposal and ProposalResponses
func InvokeOrQuery(
	spec *peer.ChaincodeSpec,
	cID string,
	txID string,
	invoke bool,
	signer msp.SigningIdentity,
	certificate tls.Certificate,
	endorserClients []*client.EndorserClient,
	deliverClients []*client.DeliverClient,
	bc common.BroadcastClient,
	cfg *Config,
) (*peer.ProposalResponse, error) {
	// Build the ChaincodeInvocationSpec message
	invocation := &peer.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error serializing identity for %s", signer.GetIdentifier()))
	}

	funcName := "invoke"
	if !invoke {
		funcName = "query"
	}

	// extract the transient field if it exists
	var tMap map[string][]byte
	if cfg.Transient != "" {
		if err := json.Unmarshal([]byte(cfg.Transient), &tMap); err != nil {
			return nil, errors.Wrap(err, "error parsing transient string")
		}
	}

	prop, txid, err := utils.CreateChaincodeProposalWithTxIDAndTransient(pcommon.HeaderType_ENDORSER_TRANSACTION, cID, invocation, creator, txID, tMap)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error creating proposal for %s", funcName))
	}

	signedProp, err := utils.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error creating signed proposal for %s", funcName))
	}
	var responses []*peer.ProposalResponse
	for _, endorser := range endorserClients {
		proposalResp, err := endorser.ProcessProposal(context.Background(), signedProp)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("error endorsing %s", funcName))
		}
		responses = append(responses, proposalResp)
	}

	if len(responses) == 0 {
		// this should only happen if some new code has introduced a bug
		return nil, errors.New("no proposal responses received - this might indicate a bug")
	}
	// all responses will be checked when the signed transaction is created.
	// for now, just set this so we check the first response's status
	proposalResp := responses[0]

	if invoke {
		if proposalResp != nil {
			if proposalResp.Response.Status >= shim.ERRORTHRESHOLD {
				return proposalResp, nil
			}
			// assemble a signed transaction (it's an Envelope message)
			env, err := utils.CreateSignedTx(prop, signer, responses...)
			if err != nil {
				return proposalResp, errors.WithMessage(err, "could not assemble transaction")
			}
			var dg *deliverGroup
			var ctx context.Context
			if cfg.WaitForEvent {
				var cancelFunc context.CancelFunc
				ctx, cancelFunc = context.WithTimeout(context.Background(), cfg.WaitForEventTimeout)
				defer cancelFunc()

				dg = newDeliverGroup(deliverClients, cfg.PeerAddresses, certificate, cfg.ChannelID, txid)
				// connect to deliver service on all peers
				err := dg.Connect(ctx)
				if err != nil {
					return nil, err
				}
			}

			// send the envelope for ordering
			if err = bc.Send(env); err != nil {
				return proposalResp, errors.WithMessage(err, fmt.Sprintf("error sending transaction for %s", funcName))
			}

			if dg != nil && ctx != nil {
				// wait for event that contains the txid from all peers
				err = dg.Wait(ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return proposalResp, nil
}
