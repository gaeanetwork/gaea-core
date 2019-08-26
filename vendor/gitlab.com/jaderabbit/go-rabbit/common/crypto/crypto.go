package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha256"
	"errors"

	"gitlab.com/jaderabbit/go-rabbit/common/crypto/sm2"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto/sm3"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto/sm4"
)

// Hash identifies a cryptographic hash function that is implemented in another
// package.
type Hash uint

const (
	// SM2 import gitlab.com/jaderabbit/go-rabbit/common/crypto/sm2
	SM2 Hash = 1 + iota
	// SM3 import gitlab.com/jaderabbit/go-rabbit/common/crypto/sm3
	SM3
	// SM4 import gitlab.com/jaderabbit/go-rabbit/common/crypto/sm4
	SM4
	// AES aes
	AES
	// DES des
	DES
	// MD5 md5
	MD5
	// SHA256 Sha256
	SHA256
	maxHash
)

// HandlerEncrypt encrypt method
type HandlerEncrypt func(plaintext, key []byte) (ciphertext []byte, err error)

// HandlerDecrypt decrypt method
type HandlerDecrypt func(plaintext, key []byte) (ciphertext []byte, err error)

// HandlerEncryptSUM encrypt method
type HandlerEncryptSUM func(plaintext []byte) (ciphertext []byte)

// hashesEncrypt is a map consisting of crypto.Hash and handlerEncrypt
var hashesEncrypt = map[Hash]HandlerEncrypt{
	SM2: SM2Encrypt,
	SM4: SM4Encrypt,
	AES: AesEncrypt,
	DES: DesEncrypt,
}

// hashesDecrypt is a map consisting of crypto.Hash and handlerDecrypt
var hashesDecrypt = map[Hash]HandlerDecrypt{
	SM2: SM2Decrypt,
	SM4: SM4Decrypt,
	AES: AesDecrypt,
	DES: DesDecrypt,
}

// hashesEncryptSUM is a map consisting of crypto.Hash and handlerEncrypt
var hashesEncryptSUM = map[Hash]HandlerEncryptSUM{
	MD5:    MD5SUM,
	SHA256: Sha256Sum,
	SM3:    SM3SUM,
}

// GetEncrypt get the encrypt function by uint
func GetEncrypt(encryptHash Hash) (HandlerEncrypt, error) {
	if encryptHash < maxHash {
		if handler, ok := hashesEncrypt[encryptHash]; ok {
			return handler, nil
		}
		return nil, errors.New("the corresponding encryption algorithm was not found")
	}

	return nil, errors.New("crypto: Size of unknown hash function")
}

// GetDecrypt get the decrypt function by uint
func GetDecrypt(decryptHash Hash) (HandlerDecrypt, error) {
	if decryptHash < maxHash {
		if handler, ok := hashesDecrypt[decryptHash]; ok {
			return handler, nil
		}
		return nil, errors.New("the corresponding encryption algorithm was not found")
	}

	return nil, errors.New("crypto: Size of unknown hash function")
}

// GetEncryptSUM get the encrypt function by uint
func GetEncryptSUM(encryptHash Hash) (HandlerEncryptSUM, error) {
	if encryptHash < maxHash {
		if handler, ok := hashesEncryptSUM[encryptHash]; ok {
			return handler, nil
		}
		return nil, errors.New("the corresponding encryption algorithm was not found")
	}

	return nil, errors.New("crypto: Size of unknown hash function")
}

// MD5SUM md5
func MD5SUM(plaintext []byte) (ciphertext []byte) {
	if len(plaintext) == 0 {
		return
	}
	sizeByte := md5.Sum(plaintext)
	ciphertext = sizeByte[:]
	return
}

// Sha256Sum sha256
func Sha256Sum(plaintext []byte) (ciphertext []byte) {
	if len(plaintext) == 0 {
		return
	}

	sizeByte := sha256.Sum256(plaintext)
	ciphertext = sizeByte[:]
	return
}

// DesEncrypt Des Encrypt by key
func DesEncrypt(plaintext, desKey []byte) (ciphertext []byte, err error) {
	if len(plaintext) == 0 {
		return
	}

	// Reference https://www.ietf.org/rfc/rfc2451.txt
	if len(desKey) < 24 {
		buffer := make([]byte, 24-len(desKey))
		desKey = append(desKey, buffer...)
	}

	block, err := des.NewTripleDESCipher(desKey)
	if err != nil {
		return nil, err
	}

	ciphertext = Encrypt(plaintext, desKey, block)
	return
}

// DesDecrypt Des Decrypt
func DesDecrypt(plaintext, desKey []byte) (ciphertext []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to decrypt by des")
		}
	}()

	if len(plaintext) == 0 {
		return
	}

	// Reference https://www.ietf.org/rfc/rfc2451.txt
	if len(desKey) < 24 {
		buffer := make([]byte, 24-len(desKey))
		desKey = append(desKey, buffer...)
	}

	block, err := des.NewTripleDESCipher(desKey)
	if err != nil {
		return nil, err
	}

	return Decrypt(plaintext, desKey, block)
}

// SM3SUM smd
func SM3SUM(plaintext []byte) (ciphertext []byte) {
	if len(plaintext) == 0 {
		return
	}

	result := sm3.Sm3Sum(plaintext)
	ciphertext = result[:]
	return
}

// SM2Encrypt sm2 Encrypt
func SM2Encrypt(plaintext, publicKey []byte) (ciphertext []byte, err error) {
	if len(plaintext) == 0 {
		return
	}

	var pub *sm2.PublicKey
	pub, err = sm2.ParseSm2PublicKey(publicKey)
	if err != nil {
		return
	}

	ciphertext, err = pub.Encrypt(plaintext)
	return
}

// SM2Decrypt sm2 Decrypt
func SM2Decrypt(plaintext []byte, privateKey []byte) (ciphertext []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to decrypt by sm2")
		}
	}()

	if len(plaintext) == 0 {
		return
	}

	var priv *sm2.PrivateKey
	priv, err = sm2.ParsePKCS8PrivateKey(privateKey, nil)
	if err != nil {
		return
	}

	ciphertext, err = priv.Decrypt(plaintext)
	return
}

// SM4Encrypt sm4 encrypt
func SM4Encrypt(plaintext, sm4Key []byte) (ciphertext []byte, err error) {
	if len(plaintext) == 0 {
		return
	}

	if len(sm4Key) < sm4.BlockSize {
		buffer := make([]byte, sm4.BlockSize-len(sm4Key))
		sm4Key = append(sm4Key, buffer...)
	}

	block, err := sm4.NewCipher(sm4Key)
	if err != nil {
		return nil, err
	}

	ciphertext = Encrypt(plaintext, sm4Key, block)
	return
}

// SM4Decrypt sm4 decrypt
func SM4Decrypt(plaintext, sm4Key []byte) (ciphertext []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to decrypt by sm4")
		}
	}()

	if len(plaintext) == 0 {
		return
	}

	if len(sm4Key) < sm4.BlockSize {
		buffer := make([]byte, sm4.BlockSize-len(sm4Key))
		sm4Key = append(sm4Key, buffer...)
	}

	block, err := sm4.NewCipher(sm4Key)
	if err != nil {
		return nil, err
	}

	return Decrypt(plaintext, sm4Key, block)
}

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

	// if the length of ase key is less then 32, add 0 to the right of key, https://www.ietf.org/rfc/rfc5649.txt
	if len(aesKey) < 32 {
		buffer := make([]byte, 32-len(aesKey))
		aesKey = append(aesKey, buffer...)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	return Decrypt(plaintext, aesKey, block)
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
