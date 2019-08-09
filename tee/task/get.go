package task

import (
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/pkg/errors"
)

// GetByID get the tee task from the blockchain by id.
func GetByID(id string) (*tee.Task, error) {
	service, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", tee.ImplementPlatform)
	}

	result, err := service.Query(ChaincodeName, []string{MethodGet, id})
	if err != nil {
		return nil, err
	}

	var d tee.Task
	if err = json.Unmarshal([]byte(result), &d); err != nil {
		return nil, err
	}

	return &d, nil
}
