package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/protos/peer"
	"gitlab.com/jaderabbit/go-rabbit/channel"
	"gitlab.com/jaderabbit/go-rabbit/i18n"
)

// Channel for view
type Channel struct {
	ID          string          `form:"-" json:"id"`
	Name        string          `form:"name" json:"name"`
	Description string          `form:"description" json:"description"`
	Peers       []*Peer         `json:"peers"`
	Orgs        []*Organization `json:"organizations"`
}

// SaveChannelInfo the channel information to syschannel
func SaveChannelInfo(c *Channel, peerAddresses []string) error {
	if len(peerAddresses) == 0 {
		err := fmt.Errorf("Error creating channel with empty peerAddresses: %v", peerAddresses)
		logger.Error(err)
		return err
	}

	if len(peerAddresses) == 1 && strings.Contains(peerAddresses[0], ",") {
		logger.Warnf("peerAddresses are comma-separated: %v", peerAddresses)
		peerAddresses = strings.Split(peerAddresses[0], ",")
	}

	logger.Infof("Create chain with peer addresses: %s", peerAddresses)
	for _, address := range peerAddresses {
		peer, err := GetPeerInfo(address)
		if err != nil {
			return err
		}
		c.Peers = append(c.Peers, peer)
	}
	if len(c.Peers) == 0 {
		err := fmt.Errorf("Error creating channel with empty peers: %v", peerAddresses)
		logger.Error(err)
		return err
	}

	org, err := GetLocalOrg()
	if err != nil {
		return err
	}
	c.Orgs = append(c.Orgs, org)
	return SaveInSysChannel(constructKey("channel", c.ID), c)
}

// GetChannelInfos gets all the channel informations from syschannel
func GetChannelInfos() ([]*Channel, error) {
	cs := make([]*Channel, 0)
	// TODO - get from org's channels or somewhere to get channel list
	channels, err := channel.List()
	if err != nil {
		return nil, err
	}

	for _, c := range channels {
		cinfo, err := GetChannelInfo(c.ChannelId)
		if err != nil {
			logger.Errorf("Error finding channel[%s] info: %s", c.ChannelId, err)
			continue
		}
		cs = append(cs, cinfo)
	}
	return cs, nil
}

// GetChannelInfo gets the channel information from syschannel. If the id is empty, it returns all channel information.
func GetChannelInfo(id string) (*Channel, error) {
	if id == SystemChannelID || id == "" {
		return GetSysChannelInfo()
	}

	value, err := GetFromSysChannel(constructKey("channel", id))
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, i18n.ChannelNotExistsErr(id)
	}

	var c Channel
	if err = json.Unmarshal([]byte(value), &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// GetSysChannelInfo get the system channel information
func GetSysChannelInfo() (*Channel, error) {
	localOrg, err := GetLocalOrg()
	if err != nil {
		return nil, err
	}

	addresses, err := GetPeersThroughGossipService(SystemChannelID)
	if err != nil {
		return nil, err
	}
	logger.Infof("All active peers address: %v", addresses)

	peers := make([]*Peer, 0)
	for _, address := range addresses {
		peerInfo, err := GetPeerInfo(address)
		if err != nil {
			logger.Error("Error getting peer info: ", err)
			peerInfo = &Peer{Name: err.Error(), Status: peer.ServerStatus_UNDEFINED.String(), Address: address}
		}

		peers = append(peers, peerInfo)
	}

	return &Channel{
		ID:          SystemChannelID,
		Name:        SystemChannelID,
		Description: SystemChannelID + " is the first chain of the blockchain system. Please rebuild the new service. Do not use it.",
		Peers:       peers,
		Orgs:        []*Organization{localOrg},
	}, nil
}
