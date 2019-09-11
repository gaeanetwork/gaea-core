package config

import (
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

	// Logger level
	LogLevel = "info"

	// Max send and receive bytes for grpc clients and servers
	MaxRecvMsgSize = 100 * 1024 * 1024
	MaxSendMsgSize = 100 * 1024 * 1024
)

const (
	// OfficialPath the default officail config path
	OfficialPath = "/etc/gaeanetwork/gaea-core/conf"
	// CurrentPath is current relative path
	CurrentPath = "./conf"

	// EnvName default config path
	EnvName        = "GAEA_CFG_PATH"
	configFileName = "gaea"
)

var (
	configDir string

	logger = flogging.MustGetLogger("Core.Config")

	gaeaViper *viper.Viper
)

// Initialize read the rabbit.yaml configuration
func Initialize() {
	gaeaViper = viper.New()
	if err := InitConfig(gaeaViper, configFileName); err != nil {
		logger.Panicf("Failed to initial %s.yaml, err: %v", configFileName, err)
	}

	readConfigConfiguration(gaeaViper)
}

func readConfigConfiguration(viper *viper.Viper) {
	// Setup core
	ListenAddr = viper.GetString("core.ListenAddr")
	GRPCAddr = viper.GetString("core.GRPCAddr")
	PProfAddr = viper.GetString("core.PProfAddr")
	ProfileEnabled = viper.GetBool("core.ProfileEnabled")
	LogLevel = viper.GetString("core.LogLevel")
}

// GetGaeaViper contains the gaea.yaml configuration
func GetGaeaViper() *viper.Viper {
	return gaeaViper
}

// GetConfigPath get config path
func GetConfigPath() string {
	return configDir
}
