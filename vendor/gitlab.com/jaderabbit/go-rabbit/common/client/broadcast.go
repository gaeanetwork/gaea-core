package client

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// BroadcastClient for orderer broadcast
type BroadcastClient struct {
	cc *grpc.ClientConn
	orderer.AtomicBroadcast_BroadcastClient
}

// Close close the grpc client connection and AtomicBroadcast_BroadcastClient send
func (c *BroadcastClient) Close() error {
	c.AtomicBroadcast_BroadcastClient.CloseSend()
	return c.cc.Close()
}

// NewBroadcastClient creates a simple instance of the BroadcastClient interface
func NewBroadcastClient() (*BroadcastClient, error) {
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

	dc, err := orderer.NewAtomicBroadcastClient(conn).Broadcast(context.TODO())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create deliver client")
	}

	return &BroadcastClient{conn, dc}, nil
}

//Send data to orderer
func (c *BroadcastClient) Send(env *common.Envelope) error {
	if err := c.AtomicBroadcast_BroadcastClient.Send(env); err != nil {
		return errors.WithMessage(err, "could not send")
	}

	return c.getAck()
}

func (c *BroadcastClient) getAck() error {
	msg, err := c.AtomicBroadcast_BroadcastClient.Recv()
	if err != nil {
		return err
	}

	if msg.Status != common.Status_SUCCESS {
		return errors.Errorf("got unexpected status: %v -- %s", msg.Status, msg.Info)
	}

	return nil
}
