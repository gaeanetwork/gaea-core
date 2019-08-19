package crypto

import (
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// PRIVFILE is the name of the private key file
	PRIVFILE = "priv.pem"
	// PUBFILE is the name of the public key file
	PUBFILE = "pub.pem"
)

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
