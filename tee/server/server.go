package server

import (
	"net"

	pb "github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services/transmission"
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
	grpcServer := grpc.NewServer()
	pb.RegisterTransmissionServer(grpcServer, transmission.NewTransmissionService())
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
