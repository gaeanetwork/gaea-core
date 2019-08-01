package tee_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/gaeanetwork/gaea-core/tee/mock"
	"github.com/stretchr/testify/assert"
)

func Test_GetNotification(t *testing.T) {
	// Service not initialized error
	factory.DeleteSmartContractService(&mock.TeeChaincodeService{})
	_, err := tee.GetNotification("data_id")
	assert.Contains(t, err.Error(), "failed to get smart contract service")

	// Query error
	queryErr := errors.New("failed to query")
	factory.InitSmartContractService(&mock.TeeChaincodeService{Error: queryErr})
	_, err = tee.GetNotification("data_id")
	assert.Equal(t, err, queryErr)

	// Result error
	invalidResult := []byte("failed to unmarshal")
	factory.InitSmartContractService(&mock.TeeChaincodeService{Result: invalidResult})
	_, err = tee.GetNotification("data_id")
	assert.Contains(t, err.Error(), "invalid character")

	// Correct invoke
	notification := &tee.Notification{ID: "??", Status: tee.Authorized}
	result, err := json.Marshal(notification)
	assert.NoError(t, err)
	factory.InitSmartContractService(&mock.TeeChaincodeService{Result: result})

	notification1, err := tee.GetNotification("??")
	assert.NoError(t, err)
	assert.Equal(t, notification, notification1)
}
