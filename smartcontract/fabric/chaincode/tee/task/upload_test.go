package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/crypto/ecc"
	"github.com/stretchr/testify/assert"
)

// test it after download data
func testUploadHookInLocal(t *testing.T, data []byte) {
	upload, exists := hooks[PubHexForTest]
	assert.True(t, exists)

	err := upload(data)
	assert.NoError(t, err)

	resultPath := filepath.Join(resultAddress, PubHexForTest+".log")
	defer os.Remove(resultPath)
	ciphertextHex, err := ioutil.ReadFile(resultPath)
	assert.NoError(t, err)
	privBytes, err := common.HexToBytes(PrivHexForTest)
	assert.NoError(t, err)
	ciphertext, err := common.HexToBytes(string(ciphertextHex))
	assert.NoError(t, err)
	plaintext, err := ecc.DecryptByECCPrivateKey(privBytes, ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, data, plaintext)
}

// test it after download data
func testUploadHookInAzure(t *testing.T, data []byte) {
	upload, exists := hooks[PubHexForTest]
	assert.True(t, exists)

	err := upload(data)
	assert.NoError(t, err)
}
