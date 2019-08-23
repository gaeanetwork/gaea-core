package chaincode

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/hyperledger/fabric/core/container/ccintf"
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