package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"gitlab.com/jaderabbit/go-rabbit/chaincode/sdk"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/system"
)

// QueryUserByName query user by userName
func QueryUserByName(userName string) (*system.User, error) {
	chaincodeServer, err := sdk.GetDefaultUserServer()
	if err != nil {
		return nil, err
	}

	arrayStr := []string{"query", "user", userName}
	userID, err := chaincodeServer.Query(arrayStr)
	if err != nil {
		return nil, err
	}

	if len(userID) == 0 {
		return nil, fmt.Errorf("user(Name:%s) does not exist", userName)
	}

	return QueryUserByID(userID)
}

// QueryUserByID query user by userID
func QueryUserByID(userID string) (*system.User, error) {
	chaincodeServer, err := sdk.GetDefaultUserServer()
	if err != nil {
		return nil, err
	}

	arrayStr := []string{"query", "", userID}
	userStr, err := chaincodeServer.Query(arrayStr)
	if err != nil {
		return nil, err
	}

	if len(userStr) == 0 {
		return nil, fmt.Errorf("user(ID:%s) does not exist", userID)
	}

	user := &system.User{}
	if err = json.Unmarshal([]byte(userStr), user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user, err:%s", err.Error())
	}
	return user, nil
}

// QueryUserByMspID query user by mspID
func QueryUserByMspID(MSPID string) (*system.User, error) {
	if len(MSPID) == 0 {
		return nil, errors.New("not specified mspid")
	}

	chaincodeServer, err := sdk.GetDefaultUserServer()
	if err != nil {
		return nil, err
	}

	arrayStr := []string{"getuser", MSPID}
	userStr, err := chaincodeServer.Query(arrayStr)
	if err != nil {
		return nil, err
	}

	if len(userStr) == 0 {
		return nil, fmt.Errorf("user(MSPID:%s) does not exist", MSPID)
	}

	user := &system.User{}
	if err = json.Unmarshal([]byte(userStr), user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user, err:%s", err.Error())
	}
	return user, nil
}
