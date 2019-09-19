package server

import (
	"net"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services/tee"
	"github.com/gaeanetwork/gaea-core/services/transmission"
	"github.com/gaeanetwork/gaea-core/services/user"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// TeeServer for tee services
type TeeServer struct {
	address string
	server  *grpc.Server
}

// NewTeeServer create a tee server by address
func NewTeeServer(address string) *TeeServer {
	//set up our server options
	var serverOpts []grpc.ServerOption

	// set max send and recv msg sizes
	serverOpts = append(serverOpts, grpc.MaxSendMsgSize(config.MaxSendMsgSize))
	serverOpts = append(serverOpts, grpc.MaxRecvMsgSize(config.MaxRecvMsgSize))

	grpcServer := grpc.NewServer(serverOpts...)
	service.RegisterTransmissionServer(grpcServer, transmission.NewTransmissionService())

	// tee data share models
	service.RegisterSharedDataServer(grpcServer, tee.NewSharedDataService())

	// user models to register and login
	user.RegisterUserServiceIntoGRPCServer(grpcServer)

	return &TeeServer{address: address, server: grpcServer}
}

// Start the tee server
func (s *TeeServer) Start() error {
	_, _, err := net.SplitHostPort(s.address)
	if err != nil {
		return errors.Wrapf(err, "address: %v is not IP:port", s.address)
	}

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Errorf("failed to listen: %v", err)
	}

	return s.server.Serve(lis)
}

// GracefulStop stops the gRPC server gracefully. It stops the server from accepting new connections and RPCs and blocks until all the pending RPCs are finished.
func (s *TeeServer) GracefulStop() {
	s.server.GracefulStop()
}

// Server return the grpc server
func (s *TeeServer) Server() *grpc.Server {
	return s.server
}
