package user

import (
	"context"
	"sync"

	pb "github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Service is used to do some shared data tasks
type Service struct {
	scService smartcontract.Service

	// TODO - SAVE IN BLOCK CHAIN OR DB
	users map[string]*pb.User

	rwMutex sync.RWMutex
}

// NewUserService create a user service
func NewUserService() *Service {
	smartcontractService, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		panic(err)
	}

	users := make(map[string]*pb.User)
	return &Service{scService: smartcontractService, users: users}
}

// RegisterUserServiceIntoGRPCServer register user service into grpc server
func RegisterUserServiceIntoGRPCServer(s *grpc.Server) {
	pb.RegisterUserServiceServer(s, NewUserService())
}

// Register a user, it will save the user public key and secret private key
func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var err error
	if err = checkRegisterRequest(req); err != nil {
		return nil, errors.Wrapf(err, "failed to check register request, req: %v", req)
	}

	s.rwMutex.RLock()
	_, exists := s.users[req.UserName]
	s.rwMutex.RUnlock()
	if exists {
		return nil, errors.Errorf("User already exists, username: %s", req.UserName)
	}

	user := &pb.User{
		Id:            uuid.New().String(),
		UserName:      req.UserName,
		Password:      req.Password,
		PublicKey:     req.PublicKey,
		SecretPrivKey: req.SecretPrivKey,
	}

	// TODO - SAVE IN BLOCK CHAIN OR DB
	s.rwMutex.Lock()
	s.users[req.UserName] = user
	s.rwMutex.Unlock()

	return &pb.RegisterResponse{User: user}, nil
}

func checkRegisterRequest(req *pb.RegisterRequest) error {
	var err error
	if err = checkUsernameLen(req.UserName); err != nil {
		return err
	}

	if err = checkPasswordLen(req.Password); err != nil {
		return err
	}

	return nil
}

// Login by username and password, return user if login successful
func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var err error
	if err = checkUsernameLen(req.UserName); err != nil {
		return nil, err
	}

	if err = checkPasswordLen(req.Password); err != nil {
		return nil, err
	}

	s.rwMutex.RLock()
	user, exists := s.users[req.UserName]
	s.rwMutex.RUnlock()
	if !exists || user.Password != req.Password {
		return nil, errors.Errorf("User does not exists or password is invalid, username: %s", req.UserName)
	}

	return &pb.LoginResponse{User: user}, nil
}

// GetUserByID get the user information by user id, return user if user id exists
func (s *Service) GetUserByID(context.Context, *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	return nil, errors.New("Not implemented")
}
