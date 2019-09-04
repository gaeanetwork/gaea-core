package main

import (
	"encoding/hex"
	"testing"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/crypto"
	"github.com/gaeanetwork/gaea-core/crypto/ecc"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/stretchr/testify/assert"
)

var (
	aesKeyForTests    = []byte("123")
	plaintextForTests = []byte("hello world!")
)

func Test_decryptDataAddress(t *testing.T) {
	stub := getTeeTaskMockStub()

	// Unencrypted
	dataAddress := string(plaintextForTests)
	dataUnencryptedInfo := &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedType: tee.UnEncrypted}
	plaintext, err := decryptDataAddress(stub, dataUnencryptedInfo)
	assert.NoError(t, err)
	assert.Equal(t, dataAddress, plaintext)

	// Data Only
	dataUnencryptedInfo = &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedType: tee.DataOnly}
	plaintext1, err := decryptDataAddress(stub, dataUnencryptedInfo)
	assert.NoError(t, err)
	assert.Equal(t, dataAddress, plaintext1)

	// Encrypted
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})

	// Address Only
	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)
	ciphertext, err := crypto.AesEncrypt(plaintextForTests, aesKeyForTests)
	assert.NoError(t, err)
	dataAddressCiphertext := hex.EncodeToString(ciphertext)
	dataAddressOnlyInfo := &tee.DataInfo{DataStoreAddress: dataAddressCiphertext, EncryptedKey: encryptedKey, EncryptedType: tee.AddressOnly}
	plaintext, err = decryptDataAddress(stub, dataAddressOnlyInfo)
	assert.NoError(t, err)
	assert.Equal(t, dataAddress, plaintext)

	// All
	dataAllInfo := &tee.DataInfo{DataStoreAddress: dataAddressCiphertext, EncryptedKey: encryptedKey, EncryptedType: tee.All}
	plaintext2, err := decryptDataAddress(stub, dataAllInfo)
	assert.NoError(t, err)
	assert.Equal(t, dataAddress, plaintext2)
}

func Test_decryptDataAddress_Error(t *testing.T) {
	stub := getTeeTaskMockStub()
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})
	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)

	// Invalid data store address hex string
	dataInfo := &tee.DataInfo{DataStoreAddress: "dataAddress", EncryptedKey: encryptedKey, EncryptedType: tee.AddressOnly}
	_, err := decryptDataAddress(stub, dataInfo)
	assert.Contains(t, err.Error(), "failed to convert the ciphertext dataAddress to bytes")

	// Invalid data store address ciphertext
	dataInfo = &tee.DataInfo{DataStoreAddress: "0x1234", EncryptedKey: encryptedKey, EncryptedType: tee.AddressOnly}
	_, err = decryptDataAddress(stub, dataInfo)
	assert.Contains(t, err.Error(), "failed to decrypt ciphertext data store address by aes algorithm")

	// Invalid encryptedKey
	ciphertext, err := crypto.AesEncrypt(plaintextForTests, aesKeyForTests)
	assert.NoError(t, err)
	dataAddress := hex.EncodeToString(ciphertext)
	dataInfo = &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedKey: "encryptedKey", EncryptedType: tee.AddressOnly}
	_, err = decryptDataAddress(stub, dataInfo)
	assert.Contains(t, err.Error(), "failed to decrypt encrypted key")

	// Invalid other encryptedKey
	encryptedKey1 := getEncryptedKeyForTests(t, aesKeyForTests[:1])
	dataInfo = &tee.DataInfo{DataStoreAddress: dataAddress, EncryptedKey: encryptedKey1, EncryptedType: tee.AddressOnly}
	_, err = decryptDataAddress(stub, dataInfo)
	assert.Contains(t, err.Error(), "failed to decrypt ciphertext data store address by aes algorithm")
}

func Test_decryptEncryptedKey(t *testing.T) {
	stub := getTeeTaskMockStub()
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})

	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)
	aesKey, err := decryptEncryptedKey(stub, encryptedKey)
	assert.NoError(t, err)
	assert.Equal(t, aesKeyForTests, aesKey)
}

func Test_decryptEncryptedKey_Error(t *testing.T) {
	stub := getTeeTaskMockStub()

	// PrivHex is empty
	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)
	_, err := decryptEncryptedKey(stub, encryptedKey)
	assert.Contains(t, err.Error(), "failed to convert the private key hex string to bytes: Error hex string is empty")

	// Invalid encrypted key hex string
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})
	_, err = decryptEncryptedKey(stub, "encryptedKey")
	assert.Contains(t, err.Error(), "faild to convert the encryptedKey to bytes")

	// Invalid encrypted key
	_, err = decryptEncryptedKey(stub, "0x1234")
	assert.Contains(t, err.Error(), "ecies: invalid public key")
}

func Test_decryptData_DataUnEncrypted(t *testing.T) {
	stub := getTeeTaskMockStub()

	dataUnencryptedInfo := &tee.DataInfo{EncryptedType: tee.UnEncrypted}
	plaintext, err := decryptData(stub, plaintextForTests, dataUnencryptedInfo)
	assert.NoError(t, err)
	assert.Equal(t, plaintextForTests, plaintext)

	// Address Only
	dataUnencryptedInfo = &tee.DataInfo{EncryptedType: tee.AddressOnly}
	plaintext1, err := decryptData(stub, plaintextForTests, dataUnencryptedInfo)
	assert.NoError(t, err)
	assert.Equal(t, plaintextForTests, plaintext1)
}

func Test_decryptData_DataEncrypted(t *testing.T) {
	stub := getTeeTaskMockStub()
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})

	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)
	ciphertext, err := crypto.AesEncrypt(plaintextForTests, aesKeyForTests)
	assert.NoError(t, err)
	dataOnlyInfo := &tee.DataInfo{EncryptedKey: encryptedKey, EncryptedType: tee.DataOnly}
	plaintext, err := decryptData(stub, []byte(hex.EncodeToString(ciphertext)), dataOnlyInfo)
	assert.NoError(t, err)
	assert.Equal(t, plaintextForTests, plaintext)

	dataAllInfo := &tee.DataInfo{EncryptedKey: encryptedKey, EncryptedType: tee.All}
	plaintext1, err := decryptData(stub, []byte(hex.EncodeToString(ciphertext)), dataAllInfo)
	assert.NoError(t, err)
	assert.Equal(t, plaintextForTests, plaintext1)
}

func Test_decryptData_DataEncrypted_Error(t *testing.T) {
	stub := getTeeTaskMockStub()
	stub.MockInit("1", [][]byte{[]byte(PrivHexForTest), []byte(PubHexForTest)})
	ciphertext, err := crypto.AesEncrypt(plaintextForTests, aesKeyForTests)
	assert.NoError(t, err)

	// Invalid format encryptedKey
	dataInfo := &tee.DataInfo{EncryptedKey: "hex.EncodeToString(encryptedKey)", EncryptedType: tee.DataOnly}
	_, err = decryptData(stub, ciphertext, dataInfo)
	assert.Contains(t, err.Error(), "failed to decrypt encrypted key")

	// Invalid ciphertext
	encryptedKey := getEncryptedKeyForTests(t, aesKeyForTests)
	dataInfo = &tee.DataInfo{EncryptedKey: encryptedKey, EncryptedType: tee.DataOnly}
	_, err = decryptData(stub, plaintextForTests, dataInfo)
	assert.Contains(t, err.Error(), "failed to convert ciphertext hex string to bytes")

	// Invalid other encryptedKey
	encryptedKey1 := getEncryptedKeyForTests(t, aesKeyForTests[:1])
	dataInfo = &tee.DataInfo{EncryptedKey: encryptedKey1, EncryptedType: tee.DataOnly}
	_, err = decryptData(stub, []byte(hex.EncodeToString(ciphertext)), dataInfo)
	assert.Contains(t, err.Error(), "failed to decrypt ciphertext data by aes algorithm")
}

func getEncryptedKeyForTests(t *testing.T, aesKey []byte) string {
	pubBytes, err := common.HexToBytes(PubHexForTest)
	assert.NoError(t, err)
	encryptedKey, err := ecc.EncryptByECCPublicKey(pubBytes, aesKey)
	assert.NoError(t, err)
	return hex.EncodeToString(encryptedKey)
}
