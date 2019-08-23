package chaincode

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/crypto/ecc"
)

// CheckArgsContainsHashAndSignatures check args for hash and signatures if they contain hashes and signatures
func CheckArgsContainsHashAndSignatures(args []string, pubkey string) error {
	length := len(args)
	if length < 2 {
		return fmt.Errorf("Invalid args structure. The args contains hash and signatures length must be greater than or equal to 2. length: %d", length)
	}

	hashIndex := length - 2
	hash, signatures := args[hashIndex], args[hashIndex+1]
	if hash1 := hex.EncodeToString(GetArgsHashBytes(args[:hashIndex])); hash != hash1 {
		return fmt.Errorf("Invalid args hash. Expecting hash: %s, Actual hash: %s, args: %v", hash1, hash, args[:hashIndex])
	}

	var sigs []string
	var err error
	if err = json.Unmarshal([]byte(signatures), &sigs); err != nil {
		return fmt.Errorf("Failed to unmarshal signatures for args[4], error: " + err.Error())
	} else if len(sigs) == 0 {
		return fmt.Errorf("Signature slice length is 0")
	}

	if err = ecc.VerifyECDSASignature(sigs[0], pubkey, hash); err != nil {
		return fmt.Errorf("Failed to verify ecdsa signature, error: " + err.Error())
	}

	return nil
}

// GetArgsHashBytes get string slice hash bytes
func GetArgsHashBytes(args []string) []byte {
	var buffer bytes.Buffer
	for _, arg := range args {
		buffer.WriteString(arg)
	}

	hash := sha256.Sum256(buffer.Bytes())
	return hash[:]
}

// GetArgsHashAndSignatures get arguments string slice hash bytes and signatures
func GetArgsHashAndSignatures(privBytes []byte, args []string) ([]byte, []byte, error) {
	hash := GetArgsHashBytes(args)
	signature, err := ecc.SignECDSA(privBytes, hash)
	if err != nil {
		return nil, nil, err
	}

	sigs := []string{signature}
	data, err := json.Marshal(sigs)
	if err != nil {
		return nil, nil, err
	}

	return hash, data, nil
}
