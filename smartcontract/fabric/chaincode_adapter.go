package fabric

import "github.com/gaeanetwork/gaea-core/smartcontract"

// Chaincode is used to send transactions to blockchain or query local ledger via chaincode
type Chaincode struct {
}

// Invoke is to send a transaction to blockchain
func (c *Chaincode) Invoke(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}

// Query is to query local ledger
func (c *Chaincode) Query(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}

// GetPlatform returns fabric
func (c *Chaincode) GetPlatform() smartcontract.Platform {
	return smartcontract.Fabric
}
