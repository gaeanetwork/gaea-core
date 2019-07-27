package mock

import "github.com/gaeanetwork/gaea-core/smartcontract"

// TeeChaincodeService for mock test
type TeeChaincodeService struct {
	Result []byte
	Error  error
}

// Invoke implement
func (t *TeeChaincodeService) Invoke(contractID string, arguments []string) (result []byte, err error) {
	if t.Error != nil {
		return nil, t.Error
	}

	return t.Result, nil
}

// Query implement
func (t *TeeChaincodeService) Query(contractID string, arguments []string) (result []byte, err error) {
	if t.Error != nil {
		return nil, t.Error
	}

	return t.Result, nil
}

// GetPlatform implement
func (t *TeeChaincodeService) GetPlatform() smartcontract.Platform { return smartcontract.Fabric }
