package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gaeanetwork/gaea-core/common"
	pb "github.com/gaeanetwork/gaea-core/protos/user"
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
