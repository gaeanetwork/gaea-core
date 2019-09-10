package user

import (
	"context"

	pb "github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"google.golang.org/grpc"
)

// Service is used to do some shared data tasks
type Service struct {
	scService smartcontract.Service

	// TODO - SAVE IN BLOCK CHAIN OR DB
	users map[string]pb.User
}

// NewUserService create a user service
func NewUserService() *Service {
	smartcontractService, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		panic(err)
	}

	users := make(map[string]pb.User)
	return &Service{scService: smartcontractService, users: users}
}

// RegisterUserServiceIntoGRPCServer register user service into grpc server
func RegisterUserServiceIntoGRPCServer(s *grpc.Server) {
	pb.RegisterUserServiceServer(s, NewUserService())
}

// Register a user, it will save the user public key and secret private key
func (s *Service) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	return nil, nil
}

// Login by username and password, return user if login successful
func (s *Service) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {

	return nil, nil
}

// GetUserByID get the user information by user id, return user if login successful
func (s *Service) GetUserByID(context.Context, *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	return nil, nil
}
