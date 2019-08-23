package client

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	defaultConnTimeout = 3 * time.Second
)

// EndorserClient for peer endorser
type EndorserClient struct {
	cc *grpc.ClientConn
	peer.EndorserClient
}

// Close close the grpc client connection
func (c *EndorserClient) Close() error {
	return c.cc.Close()
}

// GetEndorserClient returns a new endorser client. If the both the address and
// tlsRootCertFile are not provided, the target values for the client are taken
// from the configuration settings for "peer.address" and
// "peer.tls.rootcert.file"
func GetEndorserClient(address, tlsRootCertFile string) (*EndorserClient, error) {
	conn, err := createClientConn(address, tlsRootCertFile)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("Unable to create a connection to the grpc client of %s", address))
	}

	return &EndorserClient{conn, peer.NewEndorserClient(conn)}, nil
}
