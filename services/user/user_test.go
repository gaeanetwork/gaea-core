package user

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common"
	pb "github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Register(t *testing.T) {
	s := NewUserService()
	req := &pb.RegisterRequest{
		UserName: "test-username",
		Password: "test-password",
	}

	resp, err := s.Register(nil, req)
	assert.NoError(t, err)
	assert.Equal(t, resp.User.UserName, req.UserName)
	assert.Equal(t, resp.User.Password, req.Password)

	// Error - Repeated register
	resp, err = s.Register(nil, req)
	assert.Error(t, err, "User already exists")

	// Error - invalid request username
	req1 := &pb.RegisterRequest{
		Password: "test-password",
	}
	resp, err = s.Register(nil, req1)
	assert.Error(t, err, "Invalid username length")
	req1.UserName = common.GetRandomStringByLen(33)
	resp, err = s.Register(nil, req1)
	assert.Error(t, err, "Invalid username length")

	// Error - invalid request password
	req2 := &pb.RegisterRequest{
		UserName: "test-username",
	}
	resp, err = s.Register(nil, req2)
	assert.Error(t, err, "Invalid password length")
	req1.UserName = common.GetRandomStringByLen(33)
	resp, err = s.Register(nil, req1)
	assert.Error(t, err, "Invalid password length")
}

func Test_Login(t *testing.T) {
	s := NewUserService()
	req := &pb.RegisterRequest{
		UserName: "test-username",
		Password: "test-password",
	}

	resp, err := s.Register(nil, req)
	assert.NoError(t, err)

	loginReq := &pb.LoginRequest{
		UserName: resp.User.UserName,
		Password: resp.User.Password,
	}
	loginResp, err := s.Login(nil, loginReq)
	assert.NoError(t, err)
	assert.Equal(t, loginResp.User, resp.User)

	// Error - Invalid request username
	req1 := &pb.LoginRequest{
		UserName: "",
		Password: resp.User.Password,
	}
	_, err = s.Login(nil, req1)
	assert.Error(t, err, "Invalid username length")
	req1.UserName = "resp.User.UserName"
	_, err = s.Login(nil, req1)
	assert.Error(t, err, "User does not exists or password is invalid")

	// Error - Invalid request password
	req2 := &pb.LoginRequest{
		UserName: resp.User.UserName,
		Password: "",
	}
	_, err = s.Login(nil, req2)
	assert.Error(t, err, "Invalid password length")
	req2.Password = "resp.User.Password"
	_, err = s.Login(nil, req2)
	assert.Error(t, err, "User does not exists or password is invalid")
}

func Test_GetUserByID(t *testing.T) {
	s := NewUserService()

	_, err := s.GetUserByID(nil, nil)
	assert.Error(t, err, "Not implemented")
}

func Test_RegisterUserServiceIntoGRPCServer(t *testing.T) {
	assert.NotPanics(t, func() {
		grpcServer := grpc.NewServer()
		RegisterUserServiceIntoGRPCServer(grpcServer)
	})
}
