package chaincode

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/chaincode/platforms"
	"github.com/hyperledger/fabric/core/chaincode/platforms/car"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/hyperledger/fabric/core/chaincode/platforms/java"
	"github.com/hyperledger/fabric/core/chaincode/platforms/node"
	"github.com/hyperledger/fabric/core/container/ccintf"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	chainFuncName = "chaincode"
	chainCmdDes   = "Operate a chaincode: install|instantiate|invoke|package|query|signpackage|upgrade|list."
)

// Chaincode-related variables.
var (
	once      sync.Once
	rwMutex   sync.RWMutex
	mapConfig = map[string]*Config{}
	logger    = flogging.MustGetLogger("chaincodeCmd")
)

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

func (cfg *Config) newChaincodeSpec(chaincodeCtorJSON string) (*peer.ChaincodeSpec, error) {
	spec := &peer.ChaincodeSpec{}

	// Build the spec
	input := &peer.ChaincodeInput{}
	if err := proto.Unmarshal([]byte(chaincodeCtorJSON), input); err != nil {
		return spec, errors.Wrap(err, "chaincode argument error")
	}

	cfg.ChaincodeLang = strings.ToUpper(cfg.ChaincodeLang)
	spec = &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_Type(peer.ChaincodeSpec_Type_value[cfg.ChaincodeLang]),
		ChaincodeId: &peer.ChaincodeID{Path: cfg.ChaincodePath, Name: cfg.ChaincodeName, Version: cfg.ChaincodeVersion},
		Input:       input,
	}
	return spec, nil
}

// CreateChaincodeSpec create chaincode spec by config
func (cfg *Config) CreateChaincodeSpec() (*peer.ChaincodeSpec, error) {
	return &peer.ChaincodeSpec{
		Type:        peer.ChaincodeSpec_Type(peer.ChaincodeSpec_Type_value[strings.ToUpper(cfg.ChaincodeLang)]),
		ChaincodeId: &peer.ChaincodeID{Path: cfg.ChaincodePath, Name: cfg.ChaincodeName, Version: cfg.ChaincodeVersion},
		Input:       &peer.ChaincodeInput{Args: common.ConvertArrayStringToByte(cfg.ChaincodeInput)},
	}, nil
}

// ReadViperConfiguration for initialization
func ReadViperConfiguration() error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	if err := viper.UnmarshalKey("chaincodeList", &mapConfig); err != nil {
		return fmt.Errorf("Could not Unmarshal %s YAML config, err: %v", "chaincodeList", err)
	}

	logger.Debugf("Initialize the default chaincode Config:%v", mapConfig)
	return nil
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

// GetContainerName get docker container name by ccid
func GetContainerName(name, version string) string {
	ccid := ccintf.CCID{Name: name, Version: version}
	vmRegExp := regexp.MustCompile("[^a-zA-Z0-9-_.]")
	return vmRegExp.ReplaceAllString(GetPreFormatImageName(ccid), "-")
}

// GetPreFormatImageName get docker image name by ccid
func GetPreFormatImageName(ccid ccintf.CCID) string {
	name, peerID, networkID := ccid.GetName(), viper.GetString("peer.id"), viper.GetString("peer.networkId")
	if networkID != "" && peerID != "" {
		name = fmt.Sprintf("%s-%s-%s", networkID, peerID, name)
	} else if networkID != "" {
		name = fmt.Sprintf("%s-%s", networkID, name)
	} else if peerID != "" {
		name = fmt.Sprintf("%s-%s", peerID, name)
	}

	return name
}
