package ecc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"gitlab.com/jaderabbit/go-rabbit/common"
)

// ECDSASignature for ecdsa curve signature
type ECDSASignature struct {
	R, S *big.Int
}

// SignECDSA return nil only if signature is successful
func SignECDSA(privHex string, data []byte) (string, error) {
	priv, err := FromPrivHex(privHex)
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, priv, data)
	if err != nil {
		return "", fmt.Errorf("Error to sign, error: %v", err)
	}

	sig := &ECDSASignature{R: r, S: s}
	sigBytes, err := json.Marshal(sig)
	if err != nil {
		return "", err
	}

	return common.BytesToHex(sigBytes), nil
}

// VerifyECDSASignature return nil only if verification is successful
func VerifyECDSASignature(sigHex, pubHex, hashHex string) error {
	if sigHex == "" {
		return fmt.Errorf("Signature is empty")
	}

	sig, err := FromSigHex(sigHex)
	if err != nil {
		return fmt.Errorf("Failed to get ECDSASignature from sigHex: %s, error: %v", sigHex, err)
	}

	pubkey, err := FromPubHex(pubHex)
	if err != nil {
		return fmt.Errorf("Failed to get ECDSA publickey from pubHex: %s, error: %v", pubHex, err)
	}

	hash, err := common.HexToBytes(hashHex)
	if err != nil {
		return fmt.Errorf("Error getting byte slice from hashHex: %s, error: %v", hashHex, err)
	}

	if !ecdsa.Verify(pubkey, hash, sig.R, sig.S) {
		return fmt.Errorf("Failed to verify the signature, publickey: %s, hash: %s, r: %d, s: %d", pubHex, hashHex, sig.R, sig.S)
	}

	return nil
}

// FromSigHex convert sig hex string to ECDSASignature
func FromSigHex(sig string) (*ECDSASignature, error) {
	sigBytes, err := common.HexToBytes(sig)
	if err != nil {
		return nil, fmt.Errorf("Error getting byte slice from sig hex string: %s, error: %v", sig, err)
	}

	var ecdsaSig ECDSASignature
	if err = json.Unmarshal(sigBytes, &ecdsaSig); err != nil {
		return nil, fmt.Errorf("Error to unmarshal sig bytes, error: %v", err)
	}

	return &ecdsaSig, nil
}
