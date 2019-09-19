package tee

import (
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/pkg/errors"
)

// SharedData user uploads shared data information
type SharedData struct {
	ID                     string   `json:"id"`
	Ciphertext             string   `json:"data" form:"data"`
	Hash                   string   `json:"hash" form:"hash"`
	Description            string   `json:"description" form:"description"`
	Owner                  string   `json:"owner" form:"owner"`
	CreateSecondsTimestamp int64    `json:"create_seconds"`
	UploadSecondsTimestamp int64    `json:"upload_seconds"`
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

// UploadData Upload used to upload shared data for users. After the data is uploaded, once someone else
// requests to query this data, the user will be notified and can authorize or reject the request.
func UploadData(data *SharedData, hash string) (*SharedData, error) {
	service, err := factory.GetSmartContractService(ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", ImplementPlatform)
	}

	if len(data.Signatures) == 0 {
		data.Signatures = make([]string, 1)
	}

	signature, err := json.Marshal(data.Signatures)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal data signatures, signatures: %v", data.Signatures)
	}

	result, err := service.Invoke(ChaincodeName, []string{MethodUpload, data.Ciphertext, data.Hash, data.Description, data.Owner, hash, string(signature)})
	if err != nil {
		return nil, err
	}

	var d SharedData
	if err = json.Unmarshal([]byte(result), &d); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal shared data, result: %s", result)
	}

	return &d, nil
}
