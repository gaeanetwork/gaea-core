package ecc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignECDSA(t *testing.T) {
	privBytes, pubBytes, err := GenerateECP256Keypair()
	assert.NoError(t, err)

	data := []byte("123")
	sigHex, err3 := SignECDSA(privBytes, data)
	assert.NoError(t, err3)

	err4 := VerifyECDSASignature(sigHex, hex.EncodeToString(pubBytes), hex.EncodeToString(data))
	assert.NoError(t, err4)

}

func Test_VerifyECDSASignature(t *testing.T) {
	// priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// assert.NoError(t, err)
	privBytes, pubBytes, err := GenerateECP256Keypair()
	assert.NoError(t, err)
	priv, err := FromPrivBytes(privBytes)
	assert.NoError(t, err)

	hashBytes := []byte("123")
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashBytes)
	assert.NoError(t, err)

	sig := &ECDSASignature{R: r, S: s}
	sigBytes, err := json.Marshal(sig)
	assert.NoError(t, err)

	sigHex := hex.EncodeToString(sigBytes)
	pubHex := hex.EncodeToString(pubBytes)
	hashHex := hex.EncodeToString(hashBytes)
	err = VerifyECDSASignature(sigHex, pubHex, hashHex)
	assert.NoError(t, err)

	// Invalid sigHex
	sigHex1 := "0x0123"
	err1 := VerifyECDSASignature(sigHex1, pubHex, hashHex)
	assert.Contains(t, err1.Error(), "Failed to get ECDSASignature from sigHex")

	sig1 := &ECDSASignature{R: big.NewInt(1), S: big.NewInt(23)}
	sigBytes1, err2 := json.Marshal(sig1)
	assert.NoError(t, err2)

	sigHex2 := hex.EncodeToString(sigBytes1)
	err3 := VerifyECDSASignature(sigHex2, pubHex, hashHex)
	assert.Contains(t, err3.Error(), "Failed to verify the signature")

	// Invalid pubHex
	pubHex1 := "0x0123"
	err4 := VerifyECDSASignature(sigHex, pubHex1, hashHex)
	assert.Contains(t, err4.Error(), "Failed to get ECDSA publickey from pubHex")

	_, pubBytes1, err5 := GenerateECP256Keypair()
	assert.NoError(t, err5)
	err6 := VerifyECDSASignature(sigHex, hex.EncodeToString(pubBytes1), hashHex)
	assert.Contains(t, err6.Error(), "Failed to verify the signature")

	// Invalid hashHex
	hashHex1 := "0x123"
	err7 := VerifyECDSASignature(sigHex, pubHex, hashHex1)
	assert.Contains(t, err7.Error(), "Error getting byte slice from hashHex")

	hashHex2 := "0123"
	err8 := VerifyECDSASignature(sigHex, pubHex, hashHex2)
	assert.Contains(t, err8.Error(), "Failed to verify the signature")
}
