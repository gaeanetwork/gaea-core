package factory

import (
	"sync"

	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/internal/smartcontract"
	"gitlab.com/jaderabbit/go-rabbit/internal/smartcontract/fabric"
)

var (
	// SmartContract service factories
	smartContractServiceMap = make(map[smartcontract.Platform]smartcontract.Service)

	defaultService         smartcontract.Service
	defaultServiceInitOnce sync.Once
)

// GetSmartContractService get a registed smart contract service
func GetSmartContractService(platform smartcontract.Platform) (smartcontract.Service, error) {
	service, exists := smartContractServiceMap[platform]
	if !exists {
		return nil, errors.Errorf("Could not find smart contract service, no '%s' provider", platform)
	}

	return service, nil
}

// InitSmartContractService initialize a smart contract service
func InitSmartContractService(service smartcontract.Service) {
	smartContractServiceMap[service.GetPlatform()] = service
}

// GetDefaultSmartContractService get the default smart contract service implemented by the fabric
func GetDefaultSmartContractService() smartcontract.Service {
	if defaultService == nil {
		defaultServiceInitOnce.Do(func() {
			defaultService = &fabric.Chaincode{}
			smartContractServiceMap[smartcontract.Fabric] = defaultService
		})
	}

	return defaultService
}

// DeleteSmartContractService delete this smart contract service
func DeleteSmartContractService(service smartcontract.Service) {
	if service != nil {
		platform := service.GetPlatform()
		if _, exists := smartContractServiceMap[platform]; exists {
			delete(smartContractServiceMap, platform)
		}
	}
}
