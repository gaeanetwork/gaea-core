package crypto

import (
	"crypto/aes"
	"errors"
)

const (
	algorithm = "aes-256-cbc"
	keyLength = 32
)

// AesEncrypt aes encrypt
func AesEncrypt(plaintext, aesKey []byte) (ciphertext []byte, err error) {
	if len(plaintext) == 0 {
		return
	}

	aesKey = wrapAESKeyWithPadding(aesKey)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	ciphertext = CBCEncrypt(plaintext, aesKey, block)
	return
}

// wrapAESKeyWithPadding if key is not keyLength, wrap it.
// See more details, https://www.ietf.org/rfc/rfc5649.txt.
func wrapAESKeyWithPadding(aesKey []byte) []byte {
	// If the length of ase key is less then keyLength, add 0 to the right of key
	if len(aesKey) < keyLength {
		buffer := make([]byte, keyLength-len(aesKey))
		aesKey = append(aesKey, buffer...)
	} else if len(aesKey) > keyLength { // Length over keyLength, intercepted to keyLength
		aesKey = aesKey[:keyLength]
	}

	return aesKey
}

// AesDecrypt aes decrypt
func AesDecrypt(plaintext, aesKey []byte) (ciphertext []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to decrypt by aes")
		}
	}()

	if len(plaintext) == 0 {
		return
	}

	aesKey = wrapAESKeyWithPadding(aesKey)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	return Decrypt(plaintext, aesKey, block)
}
