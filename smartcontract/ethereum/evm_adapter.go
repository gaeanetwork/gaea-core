package ethereum

import "github.com/gaeanetwork/gaea-core/smartcontract"

// EVM is used to send transactions to blockchain or query local ledger via EVM
type EVM struct {
}

// Invoke is to send a transaction to blockchain
func (c *EVM) Invoke(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}

// Query is to query local ledger
func (c *EVM) Query(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}

// GetPlatform returns ethereum
func (c *EVM) GetPlatform() smartcontract.Platform {
	return smartcontract.Ethereum
}
