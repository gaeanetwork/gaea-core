package task

import (
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/pkg/errors"
)

// CreateRequest is to create a task through this request
type CreateRequest struct {
	DataIDs       []string `form:"data_id[]"`
	AlgorithmID   string   `form:"algorithm_id"`
	ResultAddress string   `form:"result_address"`
}

// Create creates a tee task by sending a create transaction to the blockchain.
func Create(req *CreateRequest) (string, error) {
	service, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get smart contract service, platform: %s", tee.ImplementPlatform)
	}

	ids, err := json.Marshal(req.DataIDs)
	if err != nil {
		return "", err
	}

	result, err := service.Invoke(ChaincodeName, []string{MethodCreate, string(ids), req.AlgorithmID, req.ResultAddress})
	if err != nil {
		return "", fmt.Errorf("Error creating tee task, req: %v, error: %v", req, err)
	}

	return string(result), nil
}
