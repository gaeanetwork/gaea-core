package tee

import (
	"encoding/json"
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/pkg/errors"
)

// Notification data view request notification received by other users
type Notification struct {
	ID                      string      `json:"id"`
	Data                    *SharedData `json:"data"`
	Requester               string      `json:"requester"`
	RequestSecondsTimestamp int64       `json:"request_seconds_timestamp"`
	Status                  AuthStatus  `json:"auth_status"`
	AuthSecondsTimestamp    int64       `json:"auth_seconds_timestamp"`
	RefusedReason           string      `json:"refused_reason"`
	DataInfo                *DataInfo   `json:"data_info"`
}

// AuthStatus status of other user authorizations
type AuthStatus int8

// AuthStatus definition
const (
	UnAuthorized AuthStatus = iota
	Authorized
	Refused
)

func (s AuthStatus) String() string {
	return fmt.Sprintf("%d", s)
}

// DataInfo claims how data is stored and encrypted.
type DataInfo struct {
	DataStoreAddress string        `json:"data_store_address"`
	DataStoreType    DataStoreType `json:"data_store_type"`
	EncryptedKey     string        `json:"encrypted_key"`
	EncryptedType    EncryptedType `json:"encrypted_type"`
}

// EncryptedType for how to use the trusted execution environment
type EncryptedType int8

// Default crypto algorithm is AES
const (
	UnEncrypted EncryptedType = iota
	AddressOnly
	DataOnly
	All
)

// DataStoreType for how to store the trusted execution result
type DataStoreType int

const (
	// Local can only be used as a single-machine consensus environment, not for multi-machine consensus.
	Local DataStoreType = iota
	// Azure currently has a Microsoft cloud account installed inside the container for uploading and downloading.
	// Currently, data storage is not considered in other Microsoft cloud accounts.
	Azure
)

// GetNotification get notification from the chain by id
func GetNotification(id string) (*Notification, error) {
	service, err := factory.GetSmartContractService(ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", ImplementPlatform)
	}

	result, err := service.Query(ChaincodeName, []string{MethodQueryDataByID, id})
	if err != nil {
		return nil, err
	}

	var d Notification
	if err = json.Unmarshal([]byte(result), &d); err != nil {
		return nil, err
	}

	return &d, nil
}
