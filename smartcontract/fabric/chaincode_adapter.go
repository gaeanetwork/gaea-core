package fabric

import (
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/pkg/errors"
)

// Chaincode is used to send transactions to blockchain or query local ledger via chaincode
type Chaincode struct {
}

// Invoke is to send a transaction to blockchain
func (c *Chaincode) Invoke(contractID string, arguments []string) (result []byte, err error) {
	server, err := chaincode.GetChaincodeServer(contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chaincode server, chaincodeName: %s", contractID)
	}

	resultStr, err := server.Invoke(arguments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to invoke chaincode, arguments: %v", arguments)
	}

	return []byte(resultStr), nil
}

// Query is to query local ledger
func (c *Chaincode) Query(contractID string, arguments []string) (result []byte, err error) {
	server, err := chaincode.GetChaincodeServer(contractID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chaincode server, chaincodeName: %s", contractID)
	}

	resultStr, err := server.Query(arguments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to invoke chaincode, arguments: %v", arguments)
	}

	return []byte(resultStr), nil
}

// GetPlatform returns fabric
func (c *Chaincode) GetPlatform() smartcontract.Platform {
	return smartcontract.Fabric
}
