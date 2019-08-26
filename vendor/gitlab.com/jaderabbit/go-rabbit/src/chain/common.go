package chain

import (
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/sdk"
)

var (
	mapSystemChaincodeServer = make(map[string]*chaincode.Server)
	mapSDKChaincodeServer    = make(map[string]*sdk.Server)

	getServerMutex    sync.Mutex
	getSDKServerMutex sync.Mutex

	logger = flogging.MustGetLogger("src.chain")
)

// GetDefaultChaincodeServer get the default chaincode server, the default chaincode and channelID set the default
func GetDefaultChaincodeServer() (*chaincode.Server, error) {
	return getSystemChaincodeServer("default", "default")
}

func getSystemChaincodeServer(chaincodeName, channelID string) (*chaincode.Server, error) {
	getServerMutex.Lock()
	defer getServerMutex.Unlock()

	key := fmt.Sprintf("%s_%s", chaincodeName, channelID)
	systemChaincodeServer, ok := mapSystemChaincodeServer[key]
	if ok {
		return systemChaincodeServer, nil
	}

	chaincodeServer, err := chaincode.GetChaincodeServer("systemchaincode")
	if err != nil {
		return nil, err
	}

	systemChaincodeServer = chaincodeServer.CopySystemChaincode(chaincodeName, channelID)
	mapSystemChaincodeServer[key] = systemChaincodeServer

	return systemChaincodeServer, nil
}

// ConstructorSystemChaincodeServer constructe a chaincode server used by system chaincode, eg:lscc
func ConstructorSystemChaincodeServer(chaincodeName, channelID string) (*sdk.Server, error) {
	getSDKServerMutex.Lock()
	defer getSDKServerMutex.Unlock()

	key := fmt.Sprintf("sdk_%s_%s", chaincodeName, channelID)
	cc, ok := mapSDKChaincodeServer[key]
	if ok {
		return cc, nil
	}

	chaincodeServer, err := sdk.GetDefaultServer()
	if err != nil {
		return nil, err
	}

	cc, err = chaincodeServer.Constructor(chaincodeName, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to constructor chaincode server, chaincodeName:%s, channelID:%s, err:%s",
			chaincodeName, channelID, err.Error())
	}
	mapSDKChaincodeServer[key] = cc

	return cc, nil
}
