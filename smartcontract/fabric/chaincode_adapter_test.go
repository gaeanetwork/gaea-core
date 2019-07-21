package fabric

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/stretchr/testify/assert"
)

func Test_SmartContractServiceImplemented(t *testing.T) {
	var chaincode interface{}
	chaincode = &Chaincode{}
	service, ok := chaincode.(smartcontract.Service)
	assert.True(t, ok)
	assert.NotNil(t, service)
	assert.Equal(t, smartcontract.Fabric, service.GetPlatform())
}
