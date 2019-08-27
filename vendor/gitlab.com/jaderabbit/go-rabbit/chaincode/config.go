package chaincode

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/chaincode/platforms"
	"github.com/hyperledger/fabric/core/chaincode/platforms/car"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/hyperledger/fabric/core/chaincode/platforms/java"
	"github.com/hyperledger/fabric/core/chaincode/platforms/node"
	"github.com/hyperledger/fabric/peer/chaincode"
	"github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
	putils "github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
)

const (
	chainFuncName = "chaincode"
	chainCmdDes   = "Operate a chaincode: install|instantiate|invoke|package|query|signpackage|upgrade|list."
)

var logger = flogging.MustGetLogger("chaincodeCmd")

// XXX This is a terrible singleton hack, however
// it simply making a latent dependency explicit.
// It should be removed along with the other package
// scoped variables
var platformRegistry = platforms.NewRegistry(
	&golang.Platform{},
	&car.Platform{},
	&java.Platform{},
	&node.Platform{},
)

// Chaincode-related variables.
var (
	once      sync.Once
	rwMutex   sync.RWMutex
	mapConfig = map[string]*Config{}
)

// Config for chaincode
type Config struct {
	ChaincodeLang         string `mapstructure:"language" yaml:"language"`
	ChaincodePath         string `mapstructure:"path" yaml:"path"`
	ChaincodeName         string `mapstructure:"name" yaml:"name"`
	ChaincodeVersion      string `mapstructure:"version" yaml:"version"`
	ChannelID             string `mapstructure:"channelID" yaml:"channelID"`
	ChaincodeQueryRaw     bool   `mapstructure:"queryRaw" yaml:"queryRaw"`
	ChaincodeQueryHex     bool   `mapstructure:"queryHex" yaml:"queryHex"`
	ChaincodeInput        []string
	ESCC                  string
	VSCC                  string
	Policy                string
	PolicyMarshalled      []byte
	CollectionsConfigFile string
	CollectionConfigBytes []byte
	ConnectionProfile     string
	PeerAddresses         []string
	TLSRootCertFiles      []string
	Transient             string
	WaitForEvent          bool
	WaitForEventTimeout   time.Duration
	CommandName           string
}

// GetDefaultConfig by default configuration
func GetDefaultConfig() *Config {
	return &Config{
		ChaincodeLang:       "golang",
		PeerAddresses:       []string{""},
		TLSRootCertFiles:    []string{""},
		WaitForEventTimeout: 30 * time.Second,
	}
}

// GetConfig get the config by chaincode name, initialize if it is nil
func GetConfig(chaincodeName string) (*Config, error) {
	if len(chaincodeName) == 0 {
		return nil, errors.New("chaincode name is empty")
	}

	rwMutex.RLock()
	defer rwMutex.RUnlock()

	cfg, ok := mapConfig[chaincodeName]
	if !ok {
		return nil, fmt.Errorf("not support chaincode name:%s", chaincodeName)
	}

	if len(cfg.PeerAddresses) == 0 {
		cfg.PeerAddresses = []string{""}
	}

	return cfg, nil
}

// ReadViperConfiguration for initialization
func ReadViperConfiguration() error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	if err := config.GetRabbitViper().UnmarshalKey("chaincodeList", &mapConfig); err != nil {
		return fmt.Errorf("Could not Unmarshal %s YAML config, err: %v", "chaincodeList", err)
	}

	logger.Debugf("Initialize the default chaincode Config:%v", mapConfig)
	return nil
}

func (cfg *Config) newChaincodeSpec(chaincodeCtorJSON string) (*pb.ChaincodeSpec, error) {
	spec := &pb.ChaincodeSpec{}

	// Build the spec
	input := &pb.ChaincodeInput{}
	if err := proto.Unmarshal([]byte(chaincodeCtorJSON), input); err != nil {
		return spec, errors.Wrap(err, "chaincode argument error")
	}

	cfg.ChaincodeLang = strings.ToUpper(cfg.ChaincodeLang)
	spec = &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[cfg.ChaincodeLang]),
		ChaincodeId: &pb.ChaincodeID{Path: cfg.ChaincodePath, Name: cfg.ChaincodeName, Version: cfg.ChaincodeVersion},
		Input:       input,
	}
	return spec, nil
}

func (cfg *Config) chaincodeInvokeOrQuery(chaincodeCtorJSON string, invoke bool, cf *chaincode.ChaincodeCmdFactory) (resp *peer.Response, err error) {
	spec, err := cfg.newChaincodeSpec(chaincodeCtorJSON)
	if err != nil {
		return nil, err
	}

	// call with empty txid to ensure production code generates a txid.
	// otherwise, tests can explicitly set their own txid
	txID := ""

	proposalResp, err := chaincode.ChaincodeInvokeOrQuery(
		spec,
		cfg.ChannelID,
		txID,
		invoke,
		cf.Signer,
		cf.Certificate,
		cf.EndorserClients,
		cf.DeliverClients,
		cf.BroadcastClient)

	if err != nil {
		return nil, errors.Errorf("%s - proposal response: %v", err, proposalResp)
	}

	if invoke {
		logger.Debugf("ESCC invoke result: %v", proposalResp)
		pRespPayload, err := putils.GetProposalResponsePayload(proposalResp.Payload)
		if err != nil {
			return nil, errors.WithMessage(err, "error while unmarshaling proposal response payload")
		}
		ca, err := putils.GetChaincodeAction(pRespPayload.Extension)
		if err != nil {
			return nil, errors.WithMessage(err, "error while unmarshaling chaincode action")
		}
		if proposalResp.Endorsement == nil {
			return nil, errors.Errorf("endorsement failure during invoke. response: %v", proposalResp.Response)
		}
		logger.Infof("Chaincode invoke successful. result: %v", ca.Response)
		return ca.Response, nil
	}

	if proposalResp == nil {
		return nil, errors.New("error during query: received nil proposal response")
	}
	if proposalResp.Endorsement == nil {
		return nil, errors.Errorf("endorsement failure during query. response: %v", proposalResp.Response)
	}

	if cfg.ChaincodeQueryRaw && cfg.ChaincodeQueryHex {
		return nil, fmt.Errorf("options --raw (-r) and --hex (-x) are not compatible")
	}
	if cfg.ChaincodeQueryRaw {
		fmt.Println(proposalResp.Response.Payload)
		return nil, nil
	}
	if cfg.ChaincodeQueryHex {
		fmt.Printf("%x\n", proposalResp.Response.Payload)
		return nil, nil
	}

	// fmt.Println(string(proposalResp.Response.Payload))
	return proposalResp.Response, nil
}

// CreateChaincodeSpec create chaincode spec by config
func (cfg *Config) CreateChaincodeSpec() (*pb.ChaincodeSpec, error) {
	return &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[strings.ToUpper(cfg.ChaincodeLang)]),
		ChaincodeId: &pb.ChaincodeID{Path: cfg.ChaincodePath, Name: cfg.ChaincodeName, Version: cfg.ChaincodeVersion},
		Input:       &pb.ChaincodeInput{Args: common.ConvertArrayStringToByte(cfg.ChaincodeInput)},
	}, nil
}
