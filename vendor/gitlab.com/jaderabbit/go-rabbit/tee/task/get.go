package task

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/jaderabbit/go-rabbit/internal/smartcontract/factory"
	"gitlab.com/jaderabbit/go-rabbit/tee"
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

// GetAll get all the tee task from the blockchain.
func GetAll() ([]*tee.Task, error) {
	service, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", tee.ImplementPlatform)
	}

	result, err := service.Query(ChaincodeName, []string{MethodGetAll})
	if err != nil {
		return nil, err
	}

	var dataList [][]byte
	if err = json.Unmarshal([]byte(result), &dataList); err != nil {
		return nil, err
	}

	dList := make([]*tee.Task, 0)
	for _, data := range dataList {
		var d tee.Task
		if err = json.Unmarshal([]byte(data), &d); err != nil {
			return nil, err
		}

		dList = append(dList, &d)
	}

	return dList, nil
}
