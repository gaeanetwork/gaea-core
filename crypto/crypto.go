package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AesEncrypt aes encrypt
func AesEncrypt(plaintext, aesKey []byte) (ciphertext []byte, err error) {
	if len(plaintext) == 0 {
		return
	}

	// if the length of ase key is less then 32, add 0 to the right of key, https://www.ietf.org/rfc/rfc5649.txt
	if len(aesKey) < 32 {
		buffer := make([]byte, 32-len(aesKey))
		aesKey = append(aesKey, buffer...)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	ciphertext = Encrypt(plaintext, aesKey, block)
	return
}

// Encrypt the plaintext to ciphertext
func Encrypt(plaintext, key []byte, block cipher.Block) (ciphertext []byte) {
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	ciphertext = make([]byte, len(plaintext))

	blockMode.CryptBlocks(ciphertext, plaintext)
	return
}

// Decrypt the ciphertext to plaintext, this function parameter is reversed.
func Decrypt(plaintext, key []byte, block cipher.Block) (ciphertext []byte, err error) {
	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	ciphertext = make([]byte, len(plaintext))

	blockMode.CryptBlocks(ciphertext, plaintext)
	return PKCS5UnPadding(ciphertext)
}

// PKCS5Padding CBC padding type
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding CBC unpadding type
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)

	unpadding := int(origData[length-1])
	for index := length - 1; index >= length-unpadding; index-- {
		if int(origData[index]) != unpadding {
			return nil, errors.New("pkcs5: incorrect password")
		}
	}

	return origData[:(length - unpadding)], nil
}
