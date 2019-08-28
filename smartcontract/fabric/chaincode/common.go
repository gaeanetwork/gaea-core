package chaincode

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/cmd"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/client"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/peer/common"
	pcommon "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"github.com/qqtou/viper"
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

type collectionConfigJSON struct {
	Name           string `json:"name"`
	Policy         string `json:"policy"`
	RequiredCount  int32  `json:"requiredPeerCount"`
	MaxPeerCount   int32  `json:"maxPeerCount"`
	BlockToLive    uint64 `json:"blockToLive"`
	MemberOnlyRead bool   `json:"memberOnlyRead"`
}

// getCollectionConfig retrieves the collection configuration
// from the supplied file; the supplied file must contain a
// json-formatted array of collectionConfigJSON elements
func getCollectionConfigFromFile(ccFile string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(ccFile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file '%s'", ccFile)
	}

	return getCollectionConfigFromBytes(fileBytes)
}

// getCollectionConfig retrieves the collection configuration
// from the supplied byte array; the byte array must contain a
// json-formatted array of collectionConfigJSON elements
func getCollectionConfigFromBytes(cconfBytes []byte) ([]byte, error) {
	cconf := &[]collectionConfigJSON{}
	err := json.Unmarshal(cconfBytes, cconf)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the collection configuration")
	}

	ccarray := make([]*pcommon.CollectionConfig, 0, len(*cconf))
	for _, cconfitem := range *cconf {
		p, err := cauthdsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("invalid policy %s", cconfitem.Policy))
		}

		cpc := &pcommon.CollectionPolicyConfig{
			Payload: &pcommon.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}

		cc := &pcommon.CollectionConfig{
			Payload: &pcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &pcommon.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: cconfitem.RequiredCount,
					MaximumPeerCount:  cconfitem.MaxPeerCount,
					BlockToLive:       cconfitem.BlockToLive,
					MemberOnlyRead:    cconfitem.MemberOnlyRead,
				},
			},
		}

		ccarray = append(ccarray, cc)
	}

	ccp := &pcommon.CollectionConfigPackage{Config: ccarray}
	return proto.Marshal(ccp)
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

// InitCmdFactory init the ChaincodeCmdFactory with default clients
func InitCmdFactory(isEndorserRequired, isOrdererRequired bool, cfg *Config) (*cmd.ChaincodeCmdFactory, error) {
	var err error
	var endorserClients []*client.EndorserClient
	var deliverClients []*client.DeliverClient
	if isEndorserRequired {
		if err = validatePeerConnectionParameters(cfg); err != nil {
			return nil, errors.WithMessage(err, "error validating peer connection parameters")
		}

		if endorserClients, deliverClients, err = getEndorserAndDeliverClients(cfg); err != nil {
			return nil, errors.WithMessage(err, "error getting endorser and deliver clients")
		}
	}

	certificate, err := common.GetCertificateFnc()
	if err != nil {
		return nil, errors.WithMessage(err, "error getting client cerificate")
	}

	signer, err := common.GetDefaultSignerFnc()
	if err != nil {
		return nil, errors.WithMessage(err, "error getting default signer")
	}

	var broadcastClient *client.BroadcastClient
	if isOrdererRequired {
		if len(common.OrderingEndpoint) == 0 {
			if len(endorserClients) == 0 {
				return nil, errors.New("orderer is required, but no ordering endpoint or endorser client supplied")
			}
			endorserClient := endorserClients[0]

			orderingEndpoints, err := common.GetOrdererEndpointOfChainFnc(cfg.ChannelID, signer, endorserClient)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("error getting channel (%s) orderer endpoint", cfg.ChannelID))
			}
			if len(orderingEndpoints) == 0 {
				return nil, errors.Errorf("no orderer endpoints retrieved for channel %s", cfg.ChannelID)
			}
			logger.Infof("Retrieved channel (%s) orderer endpoint: %s", cfg.ChannelID, orderingEndpoints[0])
			// override viper env, check value first before set it, otherwise will fatal on "concurrent map read and map write"
			if viper.Get("orderer.address") != orderingEndpoints[0] {
				viper.Set("orderer.address", orderingEndpoints[0])
			}
		}

		broadcastClient, err = client.NewBroadcastClient()
		if err != nil {
			return nil, errors.WithMessage(err, "error getting broadcast client")
		}
	}
	return &cmd.ChaincodeCmdFactory{
		EndorserClients: endorserClients,
		DeliverClients:  deliverClients,
		Signer:          signer,
		BroadcastClient: broadcastClient,
		Certificate:     certificate,
	}, nil
}

func validatePeerConnectionParameters(cfg *Config) error {
	connectionProfile := cfg.ConnectionProfile
	if connectionProfile != common.UndefinedParamValue {
		networkConfig, err := common.GetConfig(connectionProfile)
		if err != nil {
			return err
		}

		if len(networkConfig.Channels[cfg.ChannelID].Peers) != 0 {
			for peer, peerChannelConfig := range networkConfig.Channels[cfg.ChannelID].Peers {
				if peerChannelConfig.EndorsingPeer {
					peerConfig, ok := networkConfig.Peers[peer]
					if !ok {
						return errors.Errorf("peer '%s' is defined in the channel config but doesn't have associated peer config", peer)
					}
					cfg.PeerAddresses = append(cfg.PeerAddresses, peerConfig.URL)
					cfg.TLSRootCertFiles = append(cfg.TLSRootCertFiles, peerConfig.TLSCACerts.Path)
				}
			}
		}
	}

	peerAddresses, tlsRootCertFiles, cmdName := cfg.PeerAddresses, cfg.TLSRootCertFiles, cfg.CommandName
	// currently only support multiple peer addresses for invoke
	if cmdName != "invoke" && len(peerAddresses) > 1 {
		return errors.Errorf("'%s' command can only be executed against one peer. received %d", cmdName, len(peerAddresses))
	}

	if len(tlsRootCertFiles) > len(peerAddresses) {
		logger.Warningf("received more TLS root cert files (%d) than peer addresses (%d)", len(tlsRootCertFiles), len(peerAddresses))
	}

	if viper.GetBool("peer.tls.enabled") {
		if len(tlsRootCertFiles) != len(peerAddresses) {
			return errors.Errorf("number of peer addresses (%d) does not match the number of TLS root cert files (%d)", len(peerAddresses), len(tlsRootCertFiles))
		}
	} else {
		tlsRootCertFiles = nil
	}

	return nil
}

func getEndorserAndDeliverClients(cfg *Config) ([]*client.EndorserClient, []*client.DeliverClient, error) {
	var err error
	var endorserClients []*client.EndorserClient
	var deliverClients []*client.DeliverClient
	for i, address := range cfg.PeerAddresses {
		var tlsRootCertFile string
		if cfg.TLSRootCertFiles != nil {
			tlsRootCertFile = cfg.TLSRootCertFiles[i]
		}
		endorserClient, err := client.GetEndorserClient(address, tlsRootCertFile)
		if err != nil {
			return nil, nil, errors.WithMessage(err, fmt.Sprintf("error getting endorser client for %s", cfg.CommandName))
		}
		endorserClients = append(endorserClients, endorserClient)
		deliverClient, err := client.GetDeliverClient(address, tlsRootCertFile)
		if err != nil {
			return nil, nil, errors.WithMessage(err, fmt.Sprintf("error getting deliver client for %s", cfg.CommandName))
		}
		deliverClients = append(deliverClients, deliverClient)
	}
	if len(endorserClients) == 0 {
		return nil, nil, errors.New("no endorser clients retrieved - this might indicate a bug")
	}
	return endorserClients, deliverClients, err
}

// getChaincodeDeploymentSpec get chaincode deployment spec given the chaincode spec
func getChaincodeDeploymentSpec(spec *pb.ChaincodeSpec, crtPkg bool) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	if chaincode.IsDevMode() == false && crtPkg {
		var err error
		if err = checkSpec(spec); err != nil {
			return nil, err
		}

		codePackageBytes, err = container.GetChaincodePackageBytes(platformRegistry, spec)
		if err != nil {
			err = errors.WithMessage(err, "error getting chaincode package bytes")
			return nil, err
		}
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

// checkSpec to see if chaincode resides within current package capture for language.
func checkSpec(spec *pb.ChaincodeSpec) error {
	// Don't allow nil value
	if spec == nil {
		return errors.New("expected chaincode specification, nil received")
	}

	return platformRegistry.ValidateSpec(spec.CCType(), spec.Path())
}
