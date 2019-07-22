package tee

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee/mock"
	"github.com/stretchr/testify/assert"
)

func Test_GetData(t *testing.T) {
	// Service not initialized error
	factory.DeleteSmartContractService(&mock.TeeChaincodeService{})
	_, err := GetData("data_id")
	assert.Contains(t, err.Error(), "failed to get smart contract service")

	// Query error
	queryErr := errors.New("failed to query")
	factory.InitSmartContractService(&mock.TeeChaincodeService{Error: queryErr})
	_, err = GetData("data_id")
	assert.Equal(t, err, queryErr)

	// Result error
	invalidResult := []byte("failed to unmarshal")
	factory.InitSmartContractService(&mock.TeeChaincodeService{Result: invalidResult})
	_, err = GetData("data_id")
	assert.Contains(t, err.Error(), "invalid character")

	// Correct invoke
	data := &SharedData{ID: "??", Owner: "!!"}
	result, err := json.Marshal(data)
	assert.NoError(t, err)
	factory.InitSmartContractService(&mock.TeeChaincodeService{Result: result})

	data1, err := GetData("??")
	assert.NoError(t, err)
	assert.Equal(t, data, data1)
}
