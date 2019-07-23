package tee

import (
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/pkg/errors"
)

// SharedData user uploads shared data information
type SharedData struct {
	ID                     string   `json:"id"`
	Ciphertext             string   `json:"ciphertext" form:"ciphertext"`
	Hash                   string   `json:"summary" form:"summary"`
	Description            string   `json:"description" form:"description"`
	Owner                  string   `json:"owner" form:"owner"`
	CreateSecondsTimestamp int64    `json:"createSeconds"`
	UploadSecondsTimestamp int64    `json:"uploadSeconds"`
	Signatures             []string `json:"signatures" form:"signatures"`
}

// GetData get data from the chain by id
func GetData(id string) (*SharedData, error) {
	service, err := factory.GetSmartContractService(ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", ImplementPlatform)
	}

	result, err := service.Query(ChaincodeName, []string{MethodQueryDataByID, id})
	if err != nil {
		return nil, err
	}

	var d SharedData
	if err = json.Unmarshal([]byte(result), &d); err != nil {
		return nil, err
	}

	return &d, nil
}
