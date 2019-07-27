package ethereum

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/stretchr/testify/assert"
)

func Test_SmartContractServiceImplemented(t *testing.T) {
	var evm interface{}
	evm = &EVM{}
	service, ok := evm.(smartcontract.Service)
	assert.True(t, ok)
	assert.NotNil(t, service)
	assert.Equal(t, smartcontract.Ethereum, service.GetPlatform())
}
