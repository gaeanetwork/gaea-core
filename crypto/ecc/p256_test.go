package ecc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateKeyPair(t *testing.T) {
	privBytes, pubBytes, err := GenerateECP256Keypair()
	assert.NoError(t, err)

	priv, err := FromPrivBytes(privBytes)
	assert.NoError(t, err)

	pub, err := FromPubBytes(pubBytes)
	assert.NoError(t, err)
	assert.Equal(t, priv.PublicKey.X, pub.X)
	assert.Equal(t, priv.PublicKey.Y, pub.Y)
}
