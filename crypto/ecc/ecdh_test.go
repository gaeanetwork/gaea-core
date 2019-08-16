package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/gaeanetwork/gaea-core/crypto"
	"github.com/gaeanetwork/gaea-core/crypto/keyagreement"
	"github.com/stretchr/testify/assert"
)

func Test_ECDH(t *testing.T) {
	priv1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	priv2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	assert.NotEqual(t, priv1.D, priv2.D)

	ecdh := &ECDH{}
	secretKey1, err := ecdh.GenerateSharedSecret(priv1, &priv2.PublicKey)
	assert.NoError(t, err)
	secretKey2, err := ecdh.GenerateSharedSecret(priv2, &priv1.PublicKey)
	assert.NoError(t, err)
	assert.Equal(t, secretKey1, secretKey2)
}

func Test_ECDH_DifferentCurve(t *testing.T) {
	priv1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	priv3, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	assert.NoError(t, err)

	ecdh := &ECDH{}
	secretKey1, err := ecdh.GenerateSharedSecret(priv1, &priv3.PublicKey)
	assert.NoError(t, err)
	secretKey2, err := ecdh.GenerateSharedSecret(priv3, &priv1.PublicKey)
	assert.NoError(t, err)
	assert.NotEqual(t, secretKey1, secretKey2)
}

func Test_ECDH_InvalidPriv(t *testing.T) {
	priv1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// point type error
	ecdh := &ECDH{}
	_, err = ecdh.GenerateSharedSecret(*priv1, &priv1.PublicKey)
	assert.Error(t, err)

	// invalid private key error
	priv2, err := rsa.GenerateKey(rand.Reader, 32)
	assert.NoError(t, err)
	_, err = ecdh.GenerateSharedSecret(priv2, &priv1.PublicKey)
	assert.Error(t, err)
}

func Test_ECDH_InvalidPub(t *testing.T) {
	priv1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// point type error
	ecdh := &ECDH{}
	_, err = ecdh.GenerateSharedSecret(priv1, priv1.PublicKey)
	assert.Error(t, err)

	// invalid private key error
	priv2, err := rsa.GenerateKey(rand.Reader, 32)
	assert.NoError(t, err)
	_, err = ecdh.GenerateSharedSecret(priv1, &priv2.PublicKey)
	assert.Error(t, err)
}

func Test_GetAlgorithm(t *testing.T) {
	ecdh := &ECDH{}
	assert.Equal(t, keyagreement.ECDH, ecdh.GetAlgorithm())
}

var (
	privHexForTests = "307702010104207843249525ae7f43e623f5bb2b28bb8b22420e8b07d14212c12ce367e980f568a00a06082a8648ce3d030107a14403420004deb43a5bb4c34cf8db53311d4d9f95d2356b8c011349ecb04fc00b73c303bc9dc0675f4ca45a562f589b993a94129482eb9b03f259ce8982e525927c3f70fdbe"
	pubHexForTests  = "04deb43a5bb4c34cf8db53311d4d9f95d2356b8c011349ecb04fc00b73c303bc9dc0675f4ca45a562f589b993a94129482eb9b03f259ce8982e525927c3f70fdbe"
)

func Test_Android_ECDH(t *testing.T) {
	androidPubKey := "04f2ca2888417bac66b5e7bcdbcbaefe1771f45e8ac29eef23ddc84157ab16e005bd7ca457632658220a6aa721d326961e4014dae8c789030c82640bd083f3daae"
	pubBytes, err := hex.DecodeString(androidPubKey)
	if err != nil {
		t.Fatal(err)
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	pubkey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	priv, err := FromPrivHex(privHexForTests)
	if err != nil {
		t.Fatal(err)
	}

	ecdh := &ECDH{}
	secretkey, err := ecdh.GenerateSharedSecret(priv, pubkey)
	assert.NoError(t, err)
	assert.Equal(t, "726cc22f046058c4e4173f11734c2705a83b3c9f73ad48a4b36ee476dbc6f4e2", hex.EncodeToString(secretkey))

	data := []byte("Hello World!")
	ciphertext, err := crypto.AesEncrypt(data, secretkey)
	assert.NoError(t, err)
	assert.Equal(t, "65b9269169d8896ad1a5428dc8a51465", hex.EncodeToString(ciphertext))
}

func Test_SHA256(t *testing.T) {
	data := []byte("04be8ac2b0cc27d92b102b7fa25fc2d5aeb9ea5c4dfb88c74d4f8532c1ece317c8a47c6f7232f676c6c1ec46b8ab2a6687c7575b9892ae815a5f84248a946564f2")
	hash := sha256.Sum256(data)
	assert.Equal(t, "b088cf414cbab06fff85602bbc27a3e24c96a757ee29c78c48b9eaa198686a12", hex.EncodeToString(hash[:]))
}
