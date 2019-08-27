package client

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// OrdererDeliverClient for orderer deliver, holds the necessary information to connect a client
// to a peer deliver service
type OrdererDeliverClient struct {
	cc *grpc.ClientConn
	common.DeliverClient
}

// Close close the grpc client connection and AtomicBroadcast_DeliverClient send
func (c *OrdererDeliverClient) Close() error {
	c.DeliverClient.Close()
	return c.cc.Close()
}

// NewDeliverClientForOrderer creates a new DeliverClient from an OrdererClient
func NewDeliverClientForOrderer(channelID string) (*OrdererDeliverClient, error) {
	address, override, clientConfig, err := configFromEnv("orderer")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load config for OrdererClient")
	}

	gClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create GRPCClient from confclientConfigig")
	}

	conn, err := gClient.NewConnection(address, override)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("Unable to create a connection to the grpc client of %s", address))
	}

	dc, err := orderer.NewAtomicBroadcastClient(conn).Deliver(context.TODO())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create deliver client")
	}

	var tlsCertHash []byte
	// check for client certificate and create hash if present
	if certificate := gClient.Certificate().Certificate; len(certificate) > 0 {
		tlsCertHash = util.ComputeSHA256(certificate[0])
	}

	return &OrdererDeliverClient{conn, common.DeliverClient{
		Service:     dc,
		ChannelID:   channelID,
		TLSCertHash: tlsCertHash,
	}}, nil
}
