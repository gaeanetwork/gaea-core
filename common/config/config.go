package config

import (
	"os"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

// variables
var (
	// TODO - read in config
	ListenAddr     = ":12666"
	GRPCAddr       = ":12667"
	PProfAddr      = ":12668"
	ProfileEnabled = false

	// Max send and receive bytes for grpc clients and servers
	MaxRecvMsgSize = 100 * 1024 * 1024
	MaxSendMsgSize = 100 * 1024 * 1024
)

const (
	// OfficialPath the default officail config path
	OfficialPath = "/etc/gaeanetwork/gaea-core"
	// CurrentPath is current relative path
	CurrentPath = "./"

	// EnvName default config path
	EnvName        = "GAEA_CFG_PATH"
	configFileName = "gaea"
)

var (
	configDir string

	logger = flogging.MustGetLogger("Core.Config")

	gaeaViper *viper.Viper
	// DefaultConfig for gaea-core
	DefaultConfig *Config
)

// Config for gaea-core
type Config struct {
	ListenAddr     string `mapstructure:"ListenAddr" yaml:"ListenAddr"`
	GRPCAddr       string `mapstructure:"GRPCAddr" yaml:"GRPCAddr"`
	PProfAddr      string `mapstructure:"PProfAddr" yaml:"PProfAddr"`
	ProfileEnabled bool   `mapstructure:"ProfileEnabled" yaml:"ProfileEnabled"`
}

// Initialize read the rabbit.yaml configuration
func Initialize() {
	var ok bool
	configDir, ok = os.LookupEnv(EnvName)
	if !ok {
		logger.Panicf("the environment variable %s is not set", EnvName)
	}

	if !common.FileOrFolderExists(configDir) {
		logger.Panicf("path:%s not found", configDir)
	}

	gaeaViper = viper.New()
	if err := InitConfig(gaeaViper, configFileName); err != nil {
		logger.Panicf("Failed to initial %s.yaml, err: %v", configFileName, err)
	}

	readConfigConfiguration()
}

func readConfigConfiguration() {
	if err := gaeaViper.UnmarshalKey("core", &DefaultConfig); err != nil {
		logger.Errorf("Could not Unmarshal %s YAML config, err: %v", "orderingEndpoint", err)
		return
	}
}

// GetGaeaViper contains the gaea.yaml configuration
func GetGaeaViper() *viper.Viper {
	return gaeaViper
}

// GetConfigPath get config path
func GetConfigPath() string {
	return configDir
}
