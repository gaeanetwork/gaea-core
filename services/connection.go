package services

import (
	"github.com/gaeanetwork/gaea-core/common/config"
	"google.golang.org/grpc"
)

// GetGRPCConnection gets a common grpc connection
// Don't forget to close the connection after use
func GetGRPCConnection() (*grpc.ClientConn, error) {
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithInsecure())
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxCallSendMsgSize(config.MaxSendMsgSize)))

	return grpc.Dial(config.GRPCAddr, dialOpts...)
}
