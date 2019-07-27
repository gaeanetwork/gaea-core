package factory

import (
	"sync"
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/stretchr/testify/assert"
)

func Test_GetService(t *testing.T) {
	_, err := GetSmartContractService(smartcontract.Fabric)
	assert.Contains(t, err.Error(), "Could not find smart contract service")

	service := GetDefaultSmartContractService()
	assert.Equal(t, smartcontract.Fabric, service.GetPlatform())

	service1, err := GetSmartContractService(smartcontract.Fabric)
	assert.NoError(t, err)
	assert.Equal(t, service, service1)

	// recover
	defaultService = nil
	defaultServiceInitOnce = sync.Once{}
}

func Test_InitService(t *testing.T) {
	var service smartcontract.Service = &testservice{}
	assert.NotPanics(t, func() {
		InitSmartContractService(service)
	})

	service1, err := GetSmartContractService(service.GetPlatform())
	assert.NoError(t, err)
	assert.Equal(t, service, service1)

	// override the smart contract
	GetDefaultSmartContractService()
	service2, err := GetSmartContractService(service.GetPlatform())
	assert.NoError(t, err)
	assert.NotEqual(t, service, service2)

	// reoverride
	InitSmartContractService(service)
	service3, err := GetSmartContractService(service.GetPlatform())
	assert.NoError(t, err)
	assert.Equal(t, service, service3)

	// recover
	defaultService = nil
	defaultServiceInitOnce = sync.Once{}
}

func Test_DeleteSmartContractService(t *testing.T) {
	serivce := GetDefaultSmartContractService()
	assert.NotNil(t, serivce)
	assert.Equal(t, smartcontract.Fabric, serivce.GetPlatform())

	DeleteSmartContractService(defaultService)
	_, err := GetSmartContractService(smartcontract.Fabric)
	assert.Contains(t, err.Error(), "Could not find smart contract service")
}

type testservice struct{}

func (t *testservice) Invoke(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}
func (t *testservice) Query(contractID string, arguments []string) (result []byte, err error) {
	return nil, nil
}
func (t *testservice) GetPlatform() smartcontract.Platform { return smartcontract.Fabric }
