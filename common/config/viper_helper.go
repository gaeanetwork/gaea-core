package config

import (
	"os"
	"strings"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
		if strings.Contains(err.Error(), "Unsupported Config Type") {
			return errors.Errorf("Could not find config file. "+
				"Please make sure that %s is set to a path "+
				"which contains %s.yaml", EnvName, fileName)
		}

		return errors.Wrapf(err, "error when reading %s.yaml config file", fileName)
	}

	return nil
}

// InitViper performs basic initialization of our viper-based configuration layer.
// Primary thrust is to establish the paths that should be consulted to find
// the configuration we need. If v == nil, we will initialize the global
// Viper instance
func InitViper(v *viper.Viper, configName string) error {
	configDir = os.Getenv(EnvName)
	if configDir != "" {
		// If the user has overridden the path with an envvar, its the only path
		// we will consider

		if !common.FileOrFolderExists(configDir) {
			return errors.Errorf("%s %s does not exist", EnvName, configDir)
		}

		AddConfigPath(v, configDir)
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
