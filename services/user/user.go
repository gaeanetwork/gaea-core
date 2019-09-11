package user

import (
	"context"
	"sync"

	pb "github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
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
	usernameLen := len(req.UserName)
	if usernameLen <= 0 || usernameLen > 32 {
		return errors.Errorf("Invalid username length, should be (0, 32], now: %d", usernameLen)
	}

	passwordLen := len(req.Password)
	if passwordLen <= 0 || passwordLen > 32 {
		return errors.Errorf("Invalid password length, should be (0, 32], now: %d", passwordLen)
	}

	return nil
}

// Login by username and password, return user if login successful
func (s *Service) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {

	return nil, nil
}

// GetUserByID get the user information by user id, return user if login successful
func (s *Service) GetUserByID(context.Context, *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {

	return nil, nil
}
