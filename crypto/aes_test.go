package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AES_Encrypt(t *testing.T) {
	aesKey, data := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"), []byte("dmyz.org")
	ciphertext, err := AesEncrypt(data, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, "1527e127f824de210938f54bc8555cd9", hex.EncodeToString(ciphertext))

	// Key - smaller length
	aesKey = []byte("hello")
	ciphertext, err = AesEncrypt(data, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, "fdc9769dc539780fdd6d5d3b3f11151a", hex.EncodeToString(ciphertext))

	// Key - larger length
	aesKey = []byte("helloAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	ciphertext, err = AesEncrypt(data, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, "6ece0eab57eb212f4a25941f4c4c9cb3", hex.EncodeToString(ciphertext))

	// data - larger length
	data = []byte("helloAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAdmyz.org")
	ciphertext, err = AesEncrypt(data, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, "241a6575abd0be96ce6598bd3053ad9e7a8937c5c703797ab30ff5e2cb51e006759e5b2543c280155ebb59c5ab085b87", hex.EncodeToString(ciphertext))
}

func Test_paddingKey(t *testing.T) {
	key := []byte("Adf")
	assert.Len(t, wrapAESKeyWithPadding(key), keyLength)

	key = []byte("AdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfA")
	assert.Len(t, wrapAESKeyWithPadding(key), keyLength)

	key = []byte("AdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfA")
	assert.Len(t, wrapAESKeyWithPadding(key), keyLength)
}
