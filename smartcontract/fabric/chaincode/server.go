package chaincode

import (
	"sync"

	"github.com/gaeanetwork/gaea-core/common/types"
)

var mapServer map[string]*Server
var getServerMutex sync.Mutex

// Server chaincode instance
type Server struct {
	cfg *Config
}

// GetChaincodeServer to invoke/query
func GetChaincodeServer(chaincodeName string) (*Server, error) {
	getServerMutex.Lock()
	defer getServerMutex.Unlock()
	if mapServer == nil {
		mapServer = make(map[string]*Server)
	}

	chaincodeServer, ok := mapServer[chaincodeName]
	if ok {
		return chaincodeServer, nil
	}

	chaincodeServer, err := NewChaincodeServer(chaincodeName)
	if err != nil {
		return nil, err
	}

	mapServer[chaincodeName] = chaincodeServer
	return chaincodeServer, nil
}

// NewChaincodeServer returns a system server to put or get status
func NewChaincodeServer(chaincodeName string) (*Server, error) {
	cfg, err := GetConfig(chaincodeName)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg: cfg,
	}, nil
}

// Query value of key
func (s *Server) Query(arrayStr []string) (string, error) {
	s.cfg.ChaincodeInput = arrayStr
	return query(s.cfg)
}

// Invoke invoke the chaincode
func (s *Server) Invoke(arrayStr []string) (string, error) {
	s.cfg.ChaincodeInput = arrayStr
	return invoke(s.cfg)
}

// List the chaincode list of the channel
func (s *Server) List(channelID string, getInstalledChaincodes, getInstantiatedChaincodes bool) ([]*types.ChaincodeInfo, error) {
	return list(s.cfg, channelID, getInstalledChaincodes, getInstantiatedChaincodes)
}

// Install install the chaincode by the byte of the chaincode
func (s *Server) Install(ccpackfile []byte) error {
	return install(s.cfg, ccpackfile)
}

// Instantiate Instantiate the chaincode
func (s *Server) Instantiate(arrayStr []string, channelID, chaincodeName, version string) error {
	s.cfg.ChaincodeInput = arrayStr
	s.cfg.ChannelID = channelID
	s.cfg.ChaincodeName = chaincodeName
	s.cfg.ChaincodeVersion = version
	return instantiate(s.cfg)
}

// Upgrade upgrade the chaincode
func (s *Server) Upgrade(arrayStr []string, channelID, chaincodeName, version string) error {
	s.cfg.ChaincodeInput = arrayStr
	s.cfg.ChannelID = channelID
	s.cfg.ChaincodeName = chaincodeName
	s.cfg.ChaincodeVersion = version
	return upgrade(s.cfg)
}

// Package package the chaincode
func (s *Server) Package(chaincodeName, path, version string) error {
	s.cfg.ChaincodePath = path
	s.cfg.ChaincodeName = chaincodeName
	s.cfg.ChaincodeVersion = version
	return chaincodePackage(s.cfg)
}

// SignPackage package the chaincode
func (s *Server) SignPackage(chaincodeName string) error {
	s.cfg.ChaincodeName = chaincodeName
	return signpackage(s.cfg)
}

// GetChaincodeName get the chaincode name
func (s *Server) GetChaincodeName() string {
	return s.cfg.ChaincodeName
}

// GetChannelID get the channel id of chaincode server
func (s *Server) GetChannelID() string {
	return s.cfg.ChannelID
}

// CopySystemChaincode copy chaincodeServer for system chaincode by chaincode name
func (s *Server) CopySystemChaincode(chaincodeName, channelID string) *Server {
	copyCFG := *s.cfg
	copyCFG.ChaincodeName = chaincodeName
	copyCFG.ChannelID = channelID
	return &Server{
		cfg: &copyCFG,
	}
}
