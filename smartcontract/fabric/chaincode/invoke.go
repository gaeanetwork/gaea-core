package chaincode

import (
	"github.com/pkg/errors"
)

func invoke(cfg *Config) (string, error) {
	if cfg.ChannelID == "" {
		return "", errors.New("The required parameter 'channelID' is empty. Rerun the command with -C flag")
	}

	cfg.CommandName = "invoke"
	cf, err := InitCmdFactory(true, true, cfg)
	if err != nil {
		return "", err
	}
	defer cf.Close()

	return chaincodeInvokeOrQuery(true, cf, cfg)
}
