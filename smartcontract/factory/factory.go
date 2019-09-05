package factory

import (
	"sync"

	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/ethereum"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric"
	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/pkg/errors"
)

var (
	// SmartContract service factories
	smartContractServiceMap = make(map[smartcontract.Platform]smartcontract.Service)

	// Initialize all the smart contract services
	defaultServiceInitOnce sync.Once
	rwMutex                sync.RWMutex
)

// GetSmartContractService get a registed smart contract service
func GetSmartContractService(platform smartcontract.Platform) (smartcontract.Service, error) {
	defaultServiceInitOnce.Do(defaultInitialize)

	rwMutex.RLock()
	service, exists := smartContractServiceMap[platform]
	rwMutex.RUnlock()
	if !exists {
		return nil, errors.Errorf("Could not find smart contract service, no '%s' provider", platform)
	}

	return service, nil
}

func defaultInitialize() {
	smartContractServiceMap[smartcontract.Fabric] = &fabric.Chaincode{}
	smartContractServiceMap[smartcontract.Ethereum] = &ethereum.EVM{}
}

// InitSmartContractService initialize a smart contract service
func InitSmartContractService(service smartcontract.Service) {
	defaultServiceInitOnce.Do(defaultInitialize)

	common.InitCmd(nil, []string{})
	chaincode.ReadViperConfiguration()
	rwMutex.Lock()
	smartContractServiceMap[service.GetPlatform()] = service
	rwMutex.Unlock()
}

// DeleteSmartContractService delete this smart contract service
func DeleteSmartContractService(platform smartcontract.Platform) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	if _, exists := smartContractServiceMap[platform]; exists {
		delete(smartContractServiceMap, platform)
	}
}
