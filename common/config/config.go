package config

import (
	"log"

	"github.com/gaeanetwork/gaea-core/common/glog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// variables
var (
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
	OfficialPath = "/etc/gaeanetwork/gaea-core/conf"
	// CurrentPath is current relative path
	CurrentPath = "./conf"

	// EnvName default config path
	EnvName        = "GAEA_CFG_PATH"
	configFileName = "gaea"
)

var (
	gaeaViper *viper.Viper
	logger    *zap.Logger
)

func init() {
	Load()
}

// Load the gaea.yaml
func Load() {
	gaeaViper = viper.New()
	if err := InitConfig(gaeaViper, configFileName); err != nil {
		log.Printf("Failed to initial %s.yaml, err: %v\n", configFileName, err)
	}

	readConfigConfiguration(gaeaViper)
	log.Println("Configuration initialization succeeded, log level:", glog.LogLevel)

	logger = glog.MustGetLoggerWithNamed("common")
	logger.Sugar().Debugf("print gaea viper: %v", gaeaViper)
}

func readConfigConfiguration(viper *viper.Viper) {
	// Setup core
	ListenAddr = viper.GetString("core.ListenAddr")
	GRPCAddr = viper.GetString("core.GRPCAddr")
	PProfAddr = viper.GetString("core.PProfAddr")
	ProfileEnabled = viper.GetBool("core.ProfileEnabled")
	glog.LogLevel = viper.GetString("core.LogLevel")
}

// GetGaeaViper contains the gaea.yaml configuration
func GetGaeaViper() *viper.Viper {
	return gaeaViper
}
