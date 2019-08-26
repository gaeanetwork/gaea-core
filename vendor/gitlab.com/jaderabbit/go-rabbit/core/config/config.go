package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/common"
)

const (
	// OfficialPath the default officail config path
	OfficialPath = "/etc/jaderabbit/go-rabbit"
	// CurrentPath is current relative path
	CurrentPath = "./"

	// EnvName default config path
	EnvName        = "RABBIT_CFG_PATH"
	configFileName = "rabbit"

	// FabricEnvName fabric config path
	FabricEnvName = "FABRIC_CFG_PATH"
)

var (
	configDir string

	logger = flogging.MustGetLogger("Core.Config")

	rabbitViper *viper.Viper
	// DefaultConfig for go-rabbit
	DefaultConfig *Config
)

// Config for go-rabbit
type Config struct {
	RabbitDataPath  string `mapstructure:"rabbitDataPath" yaml:"rabbitDataPath"`
	IsAdmin         bool   `mapstructure:"isAdmin" yaml:"isAdmin"`
	OrdererEndpoint string `mapstructure:"ordererEndpoint" yaml:"ordererEndpoint"`
	// Note: if you use this MspConfigPath and LocalMspID, the  mspConfigPath and localMspId of core.yaml are useless
	MspConfigPath       string `mapstructure:"mspConfigPath" yaml:"mspConfigPath"`
	LocalMspID          string `mapstructure:"localMspID" yaml:"localMspID"`
	SystemChannelID     string `mapstructure:"systemChannelID" yaml:"systemChannelID"`
	CronOpen            bool   `mapstructure:"cronOpen" yaml:"cronOpen"`
	OrdererType         string `mapstructure:"ordererType" yaml:"ordererType"`
	WarnTxNumberPerHour int64  `mapstructure:"warnTxNumberPerHour" yaml:"warnTxNumberPerHour"`
	WarnLoginNumPerDay  int64  `mapstructure:"warnLoginNumPerDay" yaml:"warnLoginNumPerDay"`
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

	os.Setenv(FabricEnvName, configDir)
	logger.Infof("Setting the environment variable %s in dir: %s", FabricEnvName, configDir)

	rabbitViper = viper.New()
	if err := InitConfig(rabbitViper, configFileName); err != nil {
		logger.Panicf("Failed to initial %s.yaml, err: %v", configFileName, err)
	}

	readConfigConfiguration()
}

func readConfigConfiguration() {
	if err := rabbitViper.UnmarshalKey("configuration", &DefaultConfig); err != nil {
		logger.Errorf("Could not Unmarshal %s YAML config, err: %v", "orderingEndpoint", err)
		return
	}

	// Note: if you use this relativePath and localMspID, the  mspConfigPath and localMspId of core.yaml are useless
	if DefaultConfig.MspConfigPath != "" {
		viper.Set("peer.mspConfigPath", DefaultConfig.MspConfigPath)
		checkIsAdminAndMspConfigPathAreConsistent(DefaultConfig.IsAdmin, DefaultConfig.MspConfigPath)
	}

	if DefaultConfig.LocalMspID != "" {
		viper.Set("peer.localMspId", DefaultConfig.LocalMspID)
	}
}

func checkIsAdminAndMspConfigPathAreConsistent(isAdmin bool, mspConfigPath string) {
	if !isAdmin && strings.Contains(mspConfigPath, "Admin") {
		logger.Warnf("DefaultConfig.IsAdmin is false but the DefaultConfig.MspConfigPath contains Admin")
	}
}

// GetRabbitViper contains the rabbit.yaml configuration
func GetRabbitViper() *viper.Viper {
	return rabbitViper
}

// GetConfigPath get config path
func GetConfigPath() string {
	return configDir
}

// InitConfig initialize fileName.yaml configuration into viper
func InitConfig(v *viper.Viper, fileName string) error {
	err := InitViper(v, fileName)
	if err != nil {
		return err
	}

	// Find and read the config file, handle errors reading the config file
	if err = v.ReadInConfig(); err != nil {
		// The version of Viper we use claims the config type isn't supported when in fact the file hasn't been found
		// Display a more helpful message to avoid confusing the user.
		if strings.Contains(fmt.Sprint(err), "Unsupported Config Type") {
			return errors.New(fmt.Sprintf("Could not find config file. "+
				"Please make sure that %s is set to a path "+
				"which contains %s.yaml", EnvName, fileName))
		}

		return errors.WithMessage(err, fmt.Sprintf("error when reading %s.yaml config file", fileName))
	}

	logger.Debugf("viper keys is %v", viper.AllKeys())
	logger.Debugf("viper setting is %v", viper.AllSettings())
	return nil
}

// InitViper performs basic initialization of our viper-based configuration layer.
// Primary thrust is to establish the paths that should be consulted to find
// the configuration we need. If v == nil, we will initialize the global
// Viper instance
func InitViper(v *viper.Viper, configName string) error {
	var altPath = os.Getenv(EnvName)
	if altPath != "" {
		// If the user has overridden the path with an envvar, its the only path
		// we will consider

		if !common.FileOrFolderExists(altPath) {
			return fmt.Errorf("%s %s does not exist", EnvName, altPath)
		}

		AddConfigPath(v, altPath)
	} else {
		// If we get here, we should use the default paths in priority order:
		//
		// *) CWD
		// *) /etc/hyperledger/fabric

		// CWD
		AddConfigPath(v, CurrentPath)

		// And finally, the official path
		if common.FileOrFolderExists(OfficialPath) {
			AddConfigPath(v, OfficialPath)
		}
	}

	// Now set the configuration file.
	if v != nil {
		v.SetConfigName(configName)
	} else {
		viper.SetConfigName(configName)
	}

	return nil
}

// AddConfigPath add a path for Viper to search for the config file in.
// Can be called multiple times to define multiple search paths. If v == nil,
// we will initialize the global Viper instance
func AddConfigPath(v *viper.Viper, p string) {
	if v != nil {
		v.AddConfigPath(p)
	} else {
		viper.AddConfigPath(p)
	}
}
