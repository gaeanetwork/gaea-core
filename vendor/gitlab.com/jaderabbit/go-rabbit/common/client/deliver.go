package client

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric/peer/chaincode/api"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// DeliverClient for peer deliver, holds the necessary information to connect a client
// to a peer deliver service
type DeliverClient struct {
	cc *grpc.ClientConn
	peer.DeliverClient
}

// Close close the grpc client connection
func (c *DeliverClient) Close() error {
	return c.cc.Close()
}

// Deliver connects the client to the Deliver RPC, to implement the api.Deliver interface Deliver method
func (c *DeliverClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (api.Deliver, error) {
	return c.DeliverClient.Deliver(ctx, opts...)
}

// DeliverFiltered connects the client to the DeliverFiltered RPC, to implement the api.Deliver interface DeliverFiltered method
func (c *DeliverClient) DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (api.Deliver, error) {
	return c.DeliverClient.DeliverFiltered(ctx, opts...)
}

// GetDeliverClient returns a client for the Deliver service for peer-specific use
// cases (i.e. DeliverFiltered)
func GetDeliverClient(address, tlsRootCertFile string) (*DeliverClient, error) {
	conn, err := createClientConn(address, tlsRootCertFile)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("Unable to create a connection to the grpc client of %s", address))
	}

	return &DeliverClient{conn, peer.NewDeliverClient(conn)}, nil
}
