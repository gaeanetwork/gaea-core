package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto/ecc"
	"gitlab.com/jaderabbit/go-rabbit/tee"
	"gitlab.com/jaderabbit/go-rabbit/tee/task"
)

// decrypt data address
func decryptDataAddress(stub shim.ChaincodeStubInterface, dataInfo *tee.DataInfo) (dataAddress string, err error) {
	switch dataInfo.EncryptedType {
	case tee.UnEncrypted, tee.DataOnly:
		dataAddress = dataInfo.DataStoreAddress
	case tee.AddressOnly, tee.All:
		message, err := common.HexToBytes(dataInfo.DataStoreAddress)
		if err != nil {
			return "", errors.Wrapf(err, "failed to convert the ciphertext dataAddress to bytes, ciphertext: %s", dataInfo.DataStoreAddress)
		}

		aesKey, err := decryptEncryptedKey(stub, dataInfo.EncryptedKey)
		if err != nil {
			return "", errors.Wrapf(err, "failed to decrypt encrypted key, encryptedKey: %s", dataInfo.EncryptedKey)
		}

		dataAddressBytes, err := crypto.AesDecrypt(message, aesKey)
		if err != nil {
			return "", errors.Wrapf(err, "failed to decrypt ciphertext data store address by aes algorithm, dataStoreAddress: %s", dataInfo.DataStoreAddress)
		}

		dataAddress = string(dataAddressBytes)
	default:
		return "", fmt.Errorf("Unimplemented encryption type, crypto type: %v", dataInfo.EncryptedType)
	}

	return dataAddress, nil
}

func decryptEncryptedKey(stub shim.ChaincodeStubInterface, encryptedKey string) ([]byte, error) {
	privHexBytes, err := stub.GetState(task.KeyPrivHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task private key hex string")
	}

	prviKey, err := ecc.FromPrivHex(string(privHexBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to pares the private key hex string")
	}

	pubKey, err := ecc.FromPubHex(encryptedKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to pares the public key hex string, key: %s", encryptedKey)
	}

	ecdh := &ecc.ECDH{}
	return ecdh.GenerateSharedSecret(prviKey, pubKey)
}

// Decrypt data
func decryptData(stub shim.ChaincodeStubInterface, ciphertextHex []byte, dataInfo *tee.DataInfo) (plaintext []byte, err error) {
	switch dataInfo.EncryptedType {
	case tee.UnEncrypted, tee.AddressOnly:
		plaintext = ciphertextHex
	case tee.DataOnly, tee.All:
		aesKey, err := decryptEncryptedKey(stub, dataInfo.EncryptedKey)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decrypt encrypted key, encryptedKey: %s", dataInfo.EncryptedKey)
		}

		ciphertext, err := common.HexToBytes(string(ciphertextHex))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert ciphertext hex string to bytes, ciphertextHex: %s", string(ciphertextHex))
		}

		plaintext, err = crypto.AesDecrypt(ciphertext, aesKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt ciphertext data by aes algorithm")
		}
	}

	return plaintext, nil
}
