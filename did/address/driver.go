package address

import (
	"errors"
	"regexp"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
)

// Block Chain Names
const (
	UnKnown  = "unknown"
	Bitcoin  = "btc"
	Ethereum = "ethereum"
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

type ethereumDriver struct {
	name string
}

func (eth ethereumDriver) resolve(address string) (string, error) {
	res := common.IsHexAddress(address)
	if res == true {
		var validAddrLower = regexp.MustCompile(`^(0x)?[0-9a-f]{40}$`)
		var validAddrUpper = regexp.MustCompile(`^(0x)?[0-9A-F]{40}$`)
		resLower := validAddrLower.MatchString(address)
		resUpper := validAddrUpper.MatchString(address)
		if !resLower && !resUpper {
			if address == common.HexToAddress(address).Hex() {
				return Ethereum, nil
			}
			return "", errors.New("not a valid ethereum address")
		}
	}
	return Ethereum, nil
}
