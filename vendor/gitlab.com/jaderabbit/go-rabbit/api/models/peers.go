package models

import (
	"fmt"

	"github.com/hyperledger/fabric/core/peer"
	"github.com/hyperledger/fabric/gossip/common"
	"github.com/hyperledger/fabric/gossip/discovery"
	"github.com/hyperledger/fabric/gossip/service"
)

// GetPeersThroughGossipService returns peers with channelid through gossip service
func GetPeersThroughGossipService(channelID string) ([]string, error) {
	peerEndpoint, err := peer.GetPeerEndpoint()
	if err != nil {
		err = fmt.Errorf("Failed to get Peer Endpoint: %s", err)
		return nil, err
	}

	localAddress := peerEndpoint.Address
	if channelID == "" {
		return convertNetworkMembersToSlice(getAllNetworkMembers(), localAddress)
	}

	chainID := common.ChainID(channelID)
	networkMembers := getNetworkMembersByChainID(chainID)
	if service.GetGossipService().SelfChannelInfo(chainID) != nil {
		return convertNetworkMembersToSlice(networkMembers, localAddress)
	}
	return convertNetworkMembersToSlice(networkMembers)
}

func getAllNetworkMembers() []discovery.NetworkMember {
	return service.GetGossipService().Peers()
}

func getNetworkMembersByChainID(chainID common.ChainID) []discovery.NetworkMember {
	return service.GetGossipService().PeersOfChannel(chainID)
}

func convertNetworkMembersToSlice(networkMembers []discovery.NetworkMember, localAddress ...string) ([]string, error) {
	peers := make([]string, 0)
	if len(networkMembers) == 0 && localAddress == nil {
		return peers, nil
	}

	peers = localAddress
	for _, peer := range networkMembers {
		addr := peer.Endpoint
		if addr == "" {
			addr = peer.InternalEndpoint
		}
		peers = append(peers, addr)
	}
	return peers, nil
}
