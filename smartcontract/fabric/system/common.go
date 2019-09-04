package system

import (
	"fmt"
	"sync"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/common/flogging"
)

var (
	mapSystemChaincodeServer = make(map[string]*chaincode.Server)

	getServerMutex sync.Mutex

	logger = flogging.MustGetLogger("fabric.system")
)

func getSystemChaincodeServer(chaincodeName, channelID string) (*chaincode.Server, error) {
	getServerMutex.Lock()
	defer getServerMutex.Unlock()

	key := fmt.Sprintf("%s_%s", chaincodeName, channelID)
	systemChaincodeServer, ok := mapSystemChaincodeServer[key]
	if ok {
		return systemChaincodeServer, nil
	}

	chaincodeServer, err := chaincode.GetChaincodeServer(tee.ChaincodeName)
	if err != nil {
		return nil, err
	}

	systemChaincodeServer = chaincodeServer.CopySystemChaincode(chaincodeName, channelID)
	mapSystemChaincodeServer[key] = systemChaincodeServer

	return systemChaincodeServer, nil
}
