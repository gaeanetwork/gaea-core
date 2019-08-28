package chaincode

import (
	"errors"
)

func query(cfg *Config) (string, error) {
	if cfg.ChannelID == "" {
		return "", errors.New("The required parameter 'channelID' is empty. Rerun the command with -C flag")
	}

	cfg.CommandName = "query"
	cf, err := InitCmdFactory(true, false, cfg)
	if err != nil {
		return "", err
	}
	defer cf.Close()

	return chaincodeInvokeOrQuery(false, cf, cfg)
}
