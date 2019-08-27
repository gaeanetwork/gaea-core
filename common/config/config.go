package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// variables
var (
	// TODO - read in config
	ListenAddr     = ":12666"
	GRPCAddr       = ":12667"
	PProfAddr      = ":12668"
	ProfileEnabled = false
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
