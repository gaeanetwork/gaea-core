package channel

import (
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/peer/common"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/common/client"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
)

// The constants needed to create chain client support,
// whether node endorsement is required,
// whether consensus is needed,
// and whether it needs to be distributed to other nodes.
const (
	EndorserRequired       bool = true
	EndorserNotRequired    bool = false
	OrdererRequired        bool = true
	OrdererNotRequired     bool = false
	PeerDeliverRequired    bool = true
	PeerDeliverNotRequired bool = false
)

// Chain creation client request timeout setting, default is 10 seconds,
// Set the logger logging variable under this package to channel client support.
var (
	timeout = 10 * time.Second
	logger  = flogging.MustGetLogger("channelClientSupport")
)

// BroadcastClientFactory to get broadcast client factory
type BroadcastClientFactory func() (*client.BroadcastClient, error)

type deliverClientIntf interface {
	GetSpecifiedBlock(num uint64) (*cb.Block, error)
	GetOldestBlock() (*cb.Block, error)
	GetNewestBlock() (*cb.Block, error)
	Close() error
}

// InitializeSystemChannel Initialize the system channel to store some common
// configuration information and is a chain of business systems that all
// nodes must join. The initialized system channel name is stored in
// rabbit.yml.
// When there is no system channel in the consensus node, the system channel
// is automatically created. If the system channel already exists, the
// initial block configuration file of the system channel is automatically
// added. If you report an error during this process, you will panic directly,
// because if you do not join the system channel, some services are not
// available, this is the control that the node must join the channel.
func InitializeSystemChannel() {
	channelID := config.DefaultConfig.SystemChannelID
	err := CreateChannel(channelID, config.DefaultConfig.OrdererEndpoint)
	if err != nil {
		// TODO if the system channel is already exists, don't panic, fetch the syschannel.block and join it.
		logger.Panicf("Panic create the system channel[id: %s], error: %v", channelID, err)
	}

	if err = JoinChain(channelID); err != nil {
		logger.Panicf("Panic join the system channel[id: %s], error: %v", channelID, err)
	}
}

// ClientSupport holds the clients used by ClientSupport.
// Client support is a component used to store some of the generated client
// information, request channel names, and requester identity signatures.
// Depending on the service, some clients may be nil.
// For details, please see the code when initializing client support.
type ClientSupport struct {
	EndorserClient   *client.EndorserClient
	Signer           msp.SigningIdentity
	BroadcastClient  common.BroadcastClient
	DeliverClient    *client.OrdererDeliverClient
	BroadcastFactory BroadcastClientFactory
	ChannelID        string
}

// MakeCreateClientSupport to create channel
func MakeCreateClientSupport(channelID string) (*ClientSupport, error) {
	cs, err := makeCommonClientSupport(channelID)
	if err != nil {
		return nil, err
	}

	ordererEndpoint := viper.GetString("orderer.address")
	if len(strings.Split(ordererEndpoint, ":")) != 2 {
		return nil, errors.Errorf("viper setting orderer.address %s is not valid or missing", ordererEndpoint)
	}

	if cs.DeliverClient, err = client.NewDeliverClientForOrderer(channelID); err != nil {
		return nil, err
	}

	logger.Infof("Orderer deliver connections initialized")
	return cs, nil
}

func makeCommonClientSupport(channelID string) (*ClientSupport, error) {
	cs := &ClientSupport{ChannelID: channelID}

	var err error
	if cs.Signer, err = common.GetDefaultSignerFnc(); err != nil {
		return nil, errors.WithMessage(err, "error getting default signer")
	}

	cs.BroadcastFactory = func() (*client.BroadcastClient, error) {
		return client.NewBroadcastClient()
	}

	return cs, nil
}

// MakeGeneralClientSupport to join/list channel
func MakeGeneralClientSupport(channelID string, peerAddress string) (*ClientSupport, error) {
	cs, err := makeCommonClientSupport(channelID)
	if err != nil {
		return nil, err
	}

	if cs.EndorserClient, err = client.GetEndorserClient(peerAddress, common.UndefinedParamValue); err != nil {
		return nil, errors.WithMessage(err, "error getting endorser client for channel")
	}

	logger.Infof("Endorser connections initialized")
	return cs, nil
}
