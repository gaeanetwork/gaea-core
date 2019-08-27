package sdk

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	sdkconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
	"gitlab.com/jaderabbit/go-rabbit/core/types"
	"gitlab.com/jaderabbit/go-rabbit/database/mongodb"
)

const (
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"

	// DefaultMspUserName when no login user, user this msp to send transactions
	DefaultMspUserName = "Admin"
)

var (
	mapServer       map[string]*Server
	getServerMutex  sync.Mutex
	orgName         string
	configOpt       core.ConfigProvider
	sdkViper        *viper.Viper
	mongodbTxResult *mongodb.ChaincodeInvokeResultConnection

	logger         = flogging.MustGetLogger("chaincode.sdk")
	ccInvokeHander = make(chan *chaincodeInvokeHander, 100)
	requestChan    = make(chan struct{}, 100)
)

// Server chaincode instance
type Server struct {
	cfg    *chaincode.Config
	client *channel.Client
}

type chaincodeInvokeHander struct {
	chaincodeInvokeResult *types.ChaincodeInvokeResult
	chaincodeServer       *Server
}

// InitSDK init function, assign values to orgName and configOpt
func InitSDK() {
	sdkViper = viper.New()
	if err := config.InitConfig(sdkViper, "sdk"); err != nil {
		logger.Panicf("Failed to initial %s, err: %v", "sdk.yaml", err)
	}

	if err := sdkViper.UnmarshalKey("client.organization", &orgName); err != nil {
		logger.Panicf("Could not Unmarshal %s YAML config, err: %v", "client.organization", err)
	}

	configPath := filepath.Join(config.GetConfigPath(), "sdk.yaml")
	configOpt = sdkconfig.FromFile(configPath)

	var err error
	mongodbTxResult, err = mongodb.GetChaincodeInvokeResultConnection()
	if err != nil {
		logger.Panicf("failed to get chaincode invoke result connection from mongodb, err:%s", err.Error())
	}

	go func() {
		for {
			select {
			case handler := <-ccInvokeHander:
				go invokeAsyncHandler(handler, false)
			default:
			}
		}
	}()
}

// GetSDKServer to invoke/query
func GetSDKServer(chaincodeName, userName string) (*Server, error) {
	getServerMutex.Lock()
	defer getServerMutex.Unlock()
	if mapServer == nil {
		mapServer = make(map[string]*Server)
	}

	mapKey := fmt.Sprintf("%s_%s", chaincodeName, userName)

	chaincodeServer, ok := mapServer[mapKey]
	if ok {
		return chaincodeServer, nil
	}

	chaincodeServer, err := NewSDKServer(chaincodeName, userName)
	if err != nil {
		return nil, err
	}

	mapServer[mapKey] = chaincodeServer
	return chaincodeServer, nil
}

// NewSDKServer returns a system server to put or get status
func NewSDKServer(chaincodeName, userName string) (*Server, error) {
	cfg, err := chaincode.GetConfig(chaincodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get chaincode config, chaincodeName:%s, err:%s", chaincodeName, err.Error())
	}

	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to new fabsdk, err:%s", err.Error())
	}

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(cfg.ChannelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("failed to new channel client, err:%s", err.Error())
	}

	return &Server{
		cfg:    cfg,
		client: client,
	}, nil
}

// Constructor constructor a server by chaincode Name and channelID used by system chaincode
func (s *Server) Constructor(chaincodeName, channelID string) (*Server, error) {
	copyCFG := *s.cfg
	copyCFG.ChaincodeName = chaincodeName
	copyCFG.ChannelID = channelID

	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to new fabsdk, err:%s", err.Error())
	}

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(copyCFG.ChannelID, fabsdk.WithUser(DefaultMspUserName), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("failed to new channel client, err:%s", err.Error())
	}

	return &Server{
		cfg:    &copyCFG,
		client: client,
	}, nil
}

// GetChaincodeName get the chaincode name
func (s *Server) GetChaincodeName() string {
	return s.cfg.ChaincodeName
}

// GetChannelID get the channel id of chaincode server
func (s *Server) GetChannelID() string {
	return s.cfg.ChannelID
}

// Query value of key
func (s *Server) Query(arrayStr []string) (string, error) {
	response, err := s.client.Query(channel.Request{ChaincodeID: s.cfg.ChaincodeName, Fcn: arrayStr[0], Args: common.ConvertArrayStringToByte(arrayStr[1:])},
		channel.WithRetry(retry.DefaultChannelOpts),
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute query method, err:%s", err.Error())
	}
	return string(response.Payload), nil
}

// Invoke invoke the chaincode
func (s *Server) Invoke(arrayStr []string, isSync ...bool) (string, error) {
	select {
	case requestChan <- struct{}{}:
	default:
		return "", fmt.Errorf("failed to invoke, too many request, chaincodeName:%s, arrayStr:%v", s.cfg.ChaincodeName, arrayStr)
	}

	tr := types.NewChaincodeInvokeResult(arrayStr)
	ccIH := newChaincodeInvokeHander(tr, s)

	if len(isSync) > 0 && isSync[0] {
		invokeAsyncHandler(ccIH, true)
		// if ErrorInfo is not empty, failed to invoke
		if len(ccIH.chaincodeInvokeResult.ErrorInfo) > 0 {
			return "", errors.New(ccIH.chaincodeInvokeResult.ErrorInfo)
		}
		return string(ccIH.chaincodeInvokeResult.Payload), nil
	}

	ccInvokeHander <- ccIH
	return tr.ID, nil
}

// BatchInvoke Execute the chain code in batches asynchronously, execute the error and return immediately,
// and return the error and IDs which are used to query the executed result.
// All execution is successful, and the IDs which are used to query the executed result are returned.
func (s *Server) BatchInvoke(arrayStrs [][]string) ([]string, error) {
	ids := []string{}
	for _, arrayStr := range arrayStrs {
		result, err := s.Invoke(arrayStr)
		if err != nil {
			return ids, err
		}
		ids = append(ids, result)
	}
	return ids, nil
}

func newChaincodeInvokeHander(tr *types.ChaincodeInvokeResult, s *Server) *chaincodeInvokeHander {
	return &chaincodeInvokeHander{chaincodeInvokeResult: tr, chaincodeServer: s}
}

func invokeAsyncHandler(ccIH *chaincodeInvokeHander, isSyncExecution bool) {
	ticker := time.Tick(5 * time.Second)
	invokeChan := make(chan int)

	go invokeHandler(ccIH, invokeChan, isSyncExecution)

	// waiting for a signal of the time or the end of execution
	select {
	case <-ticker:
		invokeTimeOutHandler(ccIH, isSyncExecution)
	case <-invokeChan:
	}

	<-requestChan
}

func invokeHandler(ccIH *chaincodeInvokeHander, invokeChan chan int, isSyncExecution bool) {
	defer func() {
		invokeChan <- 1
		close(invokeChan)
	}()

	response, err := ccIH.chaincodeServer.client.Execute(channel.Request{ChaincodeID: ccIH.chaincodeServer.cfg.ChaincodeName, Fcn: ccIH.chaincodeInvokeResult.InputArgs[0], Args: common.ConvertArrayStringToByte(ccIH.chaincodeInvokeResult.InputArgs[1:])},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		ccIH.chaincodeInvokeResult.ErrorInfo = err.Error()
	}

	ccIH.chaincodeInvokeResult.ChaincodeStatus = response.ChaincodeStatus
	ccIH.chaincodeInvokeResult.TransactionID = response.TransactionID
	ccIH.chaincodeInvokeResult.TxValidationCode = response.TxValidationCode
	ccIH.chaincodeInvokeResult.Payload = response.Payload
	ccIH.chaincodeInvokeResult.ChannelID = ccIH.chaincodeServer.cfg.ChannelID
	ccIH.chaincodeInvokeResult.ChaincodeName = ccIH.chaincodeServer.cfg.ChaincodeName

	// If the execution is synchronous, the execution result is returned directly
	// and the execution result is no longer inserted into the mongodb
	if isSyncExecution {
		return
	}

	err = mongodbTxResult.Insert(ccIH.chaincodeInvokeResult)
	if err != nil {
		logger.Errorf("failed to insert chaincode invoke result, chaincodeInvokeResult:%v", ccIH.chaincodeInvokeResult)
	}
}

func invokeTimeOutHandler(ccIH *chaincodeInvokeHander, isSyncExecution bool) {
	ccIH.chaincodeInvokeResult.ErrorInfo = "time out"
	ccIH.chaincodeInvokeResult.ChannelID = ccIH.chaincodeServer.cfg.ChannelID
	ccIH.chaincodeInvokeResult.ChaincodeName = ccIH.chaincodeServer.cfg.ChaincodeName

	// If the execution is synchronous, the time out result is returned directly
	// and the result is no longer inserted into the mongodb
	if isSyncExecution {
		return
	}

	err := mongodbTxResult.Insert(ccIH.chaincodeInvokeResult)
	if err != nil {
		logger.Errorf("failed to insert chaincode invoke result, chaincodeInvokeResult:%v", ccIH.chaincodeInvokeResult)
	}
}

// GetTxInvokeResult get result of the transaction invoked
func GetTxInvokeResult(ID string) (*types.ChaincodeInvokeResult, error) {
	return mongodbTxResult.Get(ID)
}
