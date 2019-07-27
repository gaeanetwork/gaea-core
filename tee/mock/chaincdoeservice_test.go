package mock

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_ChaincodeService(t *testing.T) {
	service := &TeeChaincodeService{
		Result: []byte("nil"),
		Error:  nil,
	}
	result, err := service.Invoke("", nil)
	assert.NoError(t, err)
	assert.Equal(t, result, service.Result)
	result1, err := service.Query("", nil)
	assert.NoError(t, err)
	assert.Equal(t, result1, service.Result)
	assert.Equal(t, smartcontract.Fabric, service.GetPlatform())

	// Error
	service = &TeeChaincodeService{
		Result: []byte("nil"),
		Error:  errors.New("nil"),
	}
	result, err = service.Invoke("", nil)
	assert.Error(t, err)
	result1, err = service.Query("", nil)
	assert.Error(t, err)
	assert.Equal(t, smartcontract.Fabric, service.GetPlatform())
}
