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

func Test_JSVerify(t *testing.T) {
	r := "19051712031440086694577267594196466496543509802038404869258494690838348448522"
	s := "9967509273308981730382860329937234314009896559079928593052727688012859598829"
	R, success := big.NewInt(0).SetString(r, 10)
	if !success {
		t.Fatal()
	}
	S, success := big.NewInt(0).SetString(s, 10)
	if !success {
		t.Fatal()
	}
	pubkey, err := FromPubHex("04441524e0cfbdabd031f61f162d82617899d3e6dd661594a8002864351592f2c5081bd0713b0d525dec61266fa78b7af27c4f06292512fc92c2eba77345b46ac3")
	if err != nil {
		panic(err)
	}

	msgHash := []byte{0, 1, 2, 3, 4,
		5, 6, 7, 8, 9,
		10}
	assert.True(t, ecdsa.Verify(pubkey, msgHash, R, S))

	d := "91082081033560243149400266323747403920535827661089212104563322716539035091"
	D, success := big.NewInt(0).SetString(d, 10)
	if !success {
		t.Fatal()
	}
	privkey := &ecdsa.PrivateKey{
		PublicKey: *pubkey,
		D:         D,
	}
	R, S, err = ecdsa.Sign(rand.Reader, privkey, msgHash)
	if err != nil {
		panic(err)
	}

	assert.True(t, ecdsa.Verify(pubkey, msgHash, R, S))
}
