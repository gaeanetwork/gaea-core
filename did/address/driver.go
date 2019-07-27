package address

import (
	"encoding/hex"
	"errors"
	"regexp"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	btccrypto "github.com/ethereum/go-ethereum/crypto"
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

func (btc btcDriver) createAddress() (string, error) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", err
	}

	params := &chaincfg.MainNetParams

	wif, err := btcutil.NewWIF(privateKey, params, true)
	if err != nil {
		return "", err
	}

	if !wif.IsForNet(params) {
		return "", errors.New("The WIF string is not valid for the `" + Bitcoin + "` network")
	}

	address, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), params)
	if err != nil {
		return "", nil
	}

	return address.EncodeAddress(), nil
}

func (btc btcDriver) verifySign(signatureSerialize string, publicKey string) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}

	sigBytes, err := hex.DecodeString(signatureSerialize)

	if err != nil {
		return false, err
	}

	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		return false, err
	}

	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)
	if verified == true {
		return true, nil
	}
	return false, nil
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

func (eth ethereumDriver) createAddress() (string, error) {
	privateKey, err := btccrypto.GenerateKey()
	if err != nil {
		return "", err
	}

	return btccrypto.PubkeyToAddress(privateKey.PublicKey).Hex(), nil
}
