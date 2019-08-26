package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/common/tools/cryptogen/metadata"
	cp "github.com/hyperledger/fabric/core/peer"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/spf13/viper"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
	"google.golang.org/grpc"
)

// Peer for view
type Peer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	OrgName  string `json:"orgName"`
	Status   string `json:"status"`
	Metadata string `json:"metadata"`
	Address  string `json:"address"`
}

func (p *Peer) String() string {
	return fmt.Sprintf("Peer[id=%s, name=%s, orgName=%s, status=%s, metadata=%s, address=%s]", p.ID, p.Name, p.OrgName, p.Status, p.Metadata, p.Address)
}

// PeerInfo local peer basic information
var PeerInfo *Peer

// Initialize peer information when the peer starts and save it in the system chain.
func Initialize() {
	peerEndpoint, err := cp.GetPeerEndpoint()
	if err != nil || peerEndpoint == nil {
		logger.Panicf("Error getting peer endpoint: %v", err)
	}

	for {
		if grpcProbe(peerEndpoint.Address) {
			break
		}
	}

	id := viper.GetString("peer.id")
	if id == "" {
		logger.Panicf("Error getting peer.id is empty in core.yaml")
	}

	PeerInfo = &Peer{
		ID:       id,
		Name:     id,
		OrgName:  viper.GetString("peer.localMspId"),
		Status:   peer.ServerStatus_STARTED.String(),
		Metadata: metadata.GetVersionInfo(),
		Address:  peerEndpoint.Address,
	}

	sleepSecond, mostSleepSecond := 1*time.Second, 180*time.Second
	for {
		// wait peer started
		time.Sleep(sleepSecond)
		if sleepSecond = sleepSecond << 1; sleepSecond > mostSleepSecond {
			sleepSecond = mostSleepSecond
		}

		// Check if the peer exists in syschannel
		orgPeer, err := GetPeerInfo(PeerInfo.Address)
		if err != nil {
			if _, ok := err.(i18n.PeerNotFoundErr); !ok {
				logger.Warnf("Error getting peer in SysChannel[id: %s, code: %s]. Retry after %s. error: %s", SystemChannelID, SystemChainCode, sleepSecond, err)
				continue
			}
		}

		// if the peer exists and is the same peer, skip saving peer
		if orgPeer != nil && orgPeer.String() == PeerInfo.String() {
			logger.Infof("Have the same peer information in %s, skip saving peer: %v", SystemChannelID, PeerInfo)
			return
		}

		if err = SaveInSysChannel(constructKey("peer", PeerInfo.Address), PeerInfo); err != nil {
			logger.Warnf("Error saving peer in SysChannel[id: %s, code: %s]. Retry after %s. error: %s", SystemChannelID, SystemChainCode, sleepSecond, err)
			continue
		}
		logger.Infof("Succeed saving peer information to the %s of %s. peer: %v", SystemChainCode, SystemChannelID, PeerInfo)
		return
	}
}

// GetPeerInfos gets peer informations by channelID, if channelID is empty, return all the peer informations.
func GetPeerInfos(channelID string) ([]*Peer, error) {
	channel, err := GetChannelInfo(channelID)
	if err != nil {
		return nil, err
	}

	return channel.Peers, nil
}

// GetPeerInfo get peer information by address
func GetPeerInfo(addrss string) (*Peer, error) {
	bs, err := GetFromSysChannel(constructKey("peer", addrss))
	if err != nil {
		return nil, err
	}

	if len(bs) == 0 {
		return nil, i18n.PeerNotFoundErr(addrss)
	}

	var peerInfo Peer
	if err = json.Unmarshal(bs, &peerInfo); err != nil {
		return nil, fmt.Errorf("Error unmarshal peer information: %s", err)
	}

	return &peerInfo, nil
}

// PeerProbe whether the node starts successfully every 5 seconds, and total probe times is five
func PeerProbe() bool {
	peerEndpoint, err := cp.GetPeerEndpoint()
	if err != nil || peerEndpoint == nil {
		logger.Errorf("Error getting peer endpoint: %v", err)
		return false
	}

	count := 0
	sleepSecond := 5 * time.Second
	for {
		if count > 5 {
			return false
		}

		if grpcProbe(peerEndpoint.Address) {
			return true
		}
		count++
		time.Sleep(sleepSecond)
	}
}

func grpcProbe(addr string) bool {
	c, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithInsecure())
	if err == nil {
		c.Close()
		return true
	}
	return false
}

// ContainsAddresses return true if the addresses is a subset of peers'Address, otherwise return false and unmatched addresses
func ContainsAddresses(peers []*Peer, addresses []string) (string, bool) {
	peersSet := make(map[string]struct{})
	for _, peer := range peers {
		peersSet[peer.Address] = struct{}{}
	}

	for _, address := range addresses {
		if _, ok := peersSet[address]; !ok {
			return address, false
		}
	}

	return "", true
}
