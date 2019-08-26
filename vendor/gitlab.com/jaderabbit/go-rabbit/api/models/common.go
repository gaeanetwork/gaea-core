package models

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/common/flogging"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
)

// SystemChannel
const (
	SystemChannelID  = "syschannel"
	SystemChainCode  = "system01"
	SystemEChainName = "echain"
)

var logger = flogging.MustGetLogger("api.models")

// SaveInSysChannel includes channel,peer,organization
func SaveInSysChannel(key string, v interface{}) error {
	if key == "" {
		return i18n.ChannelKeyEmptyErr(key)
	}

	server, err := chaincode.GetChaincodeServer(SystemChainCode)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = server.Invoke([]string{"put", key, string(bs)}); err != nil {
		return err
	}

	return nil
}

// GetFromSysChannel includes channel,peer,organization
func GetFromSysChannel(key string) ([]byte, error) {
	if key == "" {
		return nil, i18n.ChannelKeyEmptyErr(key)
	}

	server, err := chaincode.GetChaincodeServer(SystemChainCode)
	if err != nil {
		return nil, err
	}

	value, err := server.Query([]string{"get", key})
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

func constructKey(id, obtype string) string {
	return fmt.Sprintf("api~models~%s~%s", obtype, id)
}
