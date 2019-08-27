package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto/ecc/ecies"
)

const (
	// PRIVFILE is the name of the private key file
	PRIVFILE = "priv.pem"
	// PUBFILE is the name of the public key file
	PUBFILE = "pub.pem"
)

// GenerateECP256Keypair generate and return private and public key byte slices using curve p256
func GenerateECP256Keypair() (privBytes []byte, pubBytes []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate ecdsa key using curve p256, error: %v", err)
	}

	privBytes, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to marshal EC private key, error: %v", err)
	}

	pubBytes = elliptic.Marshal(elliptic.P256(), priv.X, priv.Y)
	return
}

// ExportPrivateKeytoPem export private key data to pem file
func ExportPrivateKeytoPem(fileName string, der []byte, encrypted bool) error {
	var block = &pem.Block{
		Bytes: der,
	}

	if encrypted {
		block.Type = "ENCRYPTED PRIVATE KEY"
	} else {
		block.Type = "PRIVATE KEY"
	}

	return ExportPemBlock(fileName, block)
}

// ExportPublicKeytoPem export public key data to pem file
func ExportPublicKeytoPem(fileName string, der []byte) error {
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}

	return ExportPemBlock(fileName, block)
}

// ExportPemBlock export block to pem file
func ExportPemBlock(fileName string, block *pem.Block) error {
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	return pem.Encode(file, block)
}

// ImportPemFile import pem file to get block bytes
func ImportPemFile(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode pem file data")
	}

	return block.Bytes, nil
}

// EncryptByECCPrivateKey encrypt the plaintext to ciphertext by ecc private key
func EncryptByECCPrivateKey(privBytes, plaintext []byte) ([]byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ParsePKCS8PrivateKey from private Bytes")
	}

	priv := privKey.(*ecdsa.PrivateKey)
	pub := ecies.ImportECDSAPublic(&priv.PublicKey)

	return ecies.Encrypt(rand.Reader, pub, plaintext, nil, nil)
}

// EncryptByECCPublicKey encrypt the plaintext to ciphertext by ecc public key
func EncryptByECCPublicKey(pubBytes, plaintext []byte) ([]byte, error) {
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	pubkey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	pub := ecies.ImportECDSAPublic(pubkey)

	return ecies.Encrypt(rand.Reader, pub, plaintext, nil, nil)
}

// DecryptByECCPrivateKey decrypt the ciphertext to plaintext by ecc private key
func DecryptByECCPrivateKey(privBytes, ciphertext []byte) ([]byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ParsePKCS8PrivateKey from private Bytes")
	}

	privkey := privKey.(*ecdsa.PrivateKey)
	priv := ecies.ImportECDSA(privkey)

	return priv.Decrypt(ciphertext, nil, nil)
}

// FromPubHex convert public key hex string to ecdsa public key
func FromPubHex(pubHex string) (*ecdsa.PublicKey, error) {
	pubBytes, err := hex.DecodeString(pubHex)
	if err != nil {
		return nil, fmt.Errorf("Error getting byte slice from pubHex: %s, error: %v", pubHex, err)
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	if x == nil || y == nil {
		return nil, fmt.Errorf("Error to unmarshal ecdsa publickey, pubBytes: %v", pubBytes)
	}

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// FromPrivHex convert private key hex string to ecdsa private key
func FromPrivHex(privHex string) (*ecdsa.PrivateKey, error) {
	privBytes, err := hex.DecodeString(privHex)
	if err != nil {
		return nil, fmt.Errorf("Error getting byte slice from privHex: %s, error: %v", privHex, err)
	}

	return FromPrivBytes(privBytes)
}

// FromPrivBytes convert private key bytes to ecdsa private key
func FromPrivBytes(privBytes []byte) (*ecdsa.PrivateKey, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse pkcs8 private key from private bytes")
	}

	return privKey.(*ecdsa.PrivateKey), nil
}
