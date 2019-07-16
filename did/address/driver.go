package address

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// Block Chain Names
const (
	UnKnown = "unknown"
	Bitcoin = "btc"
)

// Driver interface
type Driver interface {
	resolve(address string) (string, error)
}

type btcDriver struct {
	name string
}

func (btc btcDriver) resolve(address string) (string, error) {
	_, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return Bitcoin, nil
}
