package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"

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
