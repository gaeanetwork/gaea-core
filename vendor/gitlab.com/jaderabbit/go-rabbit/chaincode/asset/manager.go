package asset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/common"
)

// ManagerRights the right info of the managers
type ManagerRights struct {
	// key is the admin userName, and value is the fields
	AdminUserFields map[string][]string `bson:"admin_user_fields" json:"admin_user_fields"`

	// the k-v, key is userName, and the value is an array of the fields
	Fields map[string][]string `bson:"fields" json:"fields"`

	TxID string `bson:"tx_id" json:"tx_id,omitempty"`
}

// InputRightInfo input parameters
type InputRightInfo struct {
	AddField    []string `bson:"add_field" json:"add_field,omitempty"`
	RemoveField []string `bson:"remove_filed" json:"remove_filed,omitempty"`
	UserName    string   `bson:"user_name" json:"user_name,omitempty"`
}

// NewManagerRights when new a ManagerRights, make the map type of the ManagerRights
func NewManagerRights() *ManagerRights {
	return &ManagerRights{
		AdminUserFields: make(map[string][]string),
		Fields:          make(map[string][]string),
	}
}

func (mr *ManagerRights) managerFieldRight(stub shim.ChaincodeStubInterface, rightInfo string) error {
	if err := mr.validateOperateFieldRight(stub); err != nil {
		return err
	}

	inputRI := &InputRightInfo{}
	if err := json.Unmarshal([]byte(rightInfo), inputRI); err != nil {
		return err
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return errors.New("failed to get user by chaincode")
	}

	if !user.IsRoot {
		adminFields, ok := mr.AdminUserFields[user.UserName]
		if !ok {
			return fmt.Errorf("admin user(Name%s) has not field to manager", user.UserName)
		}

		noRightField := common.RemoveDuplicateArray(inputRI.AddField, adminFields)
		if len(noRightField) > 0 {
			return fmt.Errorf("admin user(Name%s) has not this field to add field right for user(%s), field:%v", user.UserName, inputRI.UserName, noRightField)
		}

		noRightField = common.RemoveDuplicateArray(inputRI.RemoveField, mr.Fields[inputRI.UserName])
		if len(noRightField) > 0 {
			return fmt.Errorf("user(Name%s) has not this field to delete, field:%v", inputRI.UserName, noRightField)
		}

		inputRI.AddField = common.GetRepeateElementArray(inputRI.AddField, adminFields)
		inputRI.RemoveField = common.GetRepeateElementArray(inputRI.RemoveField, adminFields)
	}

	if len(inputRI.AddField) > 0 {
		// add fields
		if _, ok := mr.Fields[inputRI.UserName]; ok {
			mr.Fields[inputRI.UserName] = common.MergeArray(mr.Fields[inputRI.UserName], inputRI.AddField)
		} else {
			mr.Fields[inputRI.UserName] = inputRI.AddField
		}
	}

	if len(inputRI.RemoveField) > 0 {
		if _, ok := mr.Fields[inputRI.UserName]; !ok {
			return fmt.Errorf("failed to remove field, there is no field of user(Name:%s)", inputRI.UserName)
		}

		mr.Fields[inputRI.UserName] = common.RemoveDuplicateArray(mr.Fields[inputRI.UserName], inputRI.RemoveField)
		if len(mr.Fields[inputRI.UserName]) == 0 {
			delete(mr.Fields, inputRI.UserName)
		}
	}

	byteMR, err := json.Marshal(mr)
	if err != nil {
		return err
	}

	return stub.PutState(rightsKey, byteMR)
}

func (mr *ManagerRights) validateOperateFieldRight(stub shim.ChaincodeStubInterface) error {
	mspID, err := chaincode.GetMSPID(stub)
	if err != nil {
		return fmt.Errorf("failed to get msp id, err:%s", err.Error())
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return err
	}

	if !user.IsRoot && !user.IsAdmin {
		return errors.New("failed to validate operate field without administrator permissions")
	}

	_, userMSPID, err := user.GetDefaultMSP()
	if err != nil {
		return err
	}

	if userMSPID != mspID {
		return fmt.Errorf("failed to validate operate filed right")
	}
	return nil
}

func (mr *ManagerRights) manager(stub shim.ChaincodeStubInterface, rightInfo string) error {
	inputRI := &InputRightInfo{}
	if err := json.Unmarshal([]byte(rightInfo), inputRI); err != nil {
		return err
	}

	if len(inputRI.UserName) == 0 {
		return errors.New("failed to manager field, the user name is empty")
	}

	user, err := getUserByUserName(stub, inputRI.UserName)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return fmt.Errorf("failed to manager field for amdin user, user(Name:%s) is not administrator", user.UserName)
	}

	if len(inputRI.AddField) > 0 {
		// add fields
		if _, ok := mr.AdminUserFields[inputRI.UserName]; ok {
			mr.AdminUserFields[inputRI.UserName] = common.MergeArray(mr.AdminUserFields[inputRI.UserName], inputRI.AddField)
		} else {
			mr.AdminUserFields[inputRI.UserName] = inputRI.AddField
		}
	}

	if len(inputRI.RemoveField) > 0 {
		if _, ok := mr.AdminUserFields[inputRI.UserName]; !ok {
			return fmt.Errorf("failed to remove field, there is no field of user(Name:%s)", inputRI.UserName)
		}

		mr.AdminUserFields[inputRI.UserName] = common.RemoveDuplicateArray(mr.AdminUserFields[inputRI.UserName], inputRI.RemoveField)
		if len(mr.AdminUserFields[inputRI.UserName]) == 0 {
			delete(mr.AdminUserFields, inputRI.UserName)
		}
	}

	byteMR, err := json.Marshal(mr)
	if err != nil {
		return err
	}

	return stub.PutState(rightsKey, byteMR)
}

func getRights(stub shim.ChaincodeStubInterface) (*ManagerRights, error) {
	byteRightKey, err := stub.GetState(rightsKey)
	if err != nil {
		return nil, err
	}

	mr := NewManagerRights()
	if len(byteRightKey) == 0 {
		return mr, nil
	}

	if err = json.Unmarshal(byteRightKey, mr); err != nil {
		return nil, err
	}
	return mr, nil
}

// Manager add or remove the user who can manager field
// the mspinfo format: ["qqtou","qqtou-zhangsan"]
func manager(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsLength(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	rightInfo := args[0]

	if err := validateOwnerRight(stub); err != nil {
		return shim.Error(fmt.Sprintf("failed to validate owner right, err:%s", err.Error()))
	}

	mr, err := getRights(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get the right, err:%s", err.Error()))
	}

	if err := mr.manager(stub, rightInfo); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// managerFieldRight manager the user that operate the field right, eg: add a query right for a user;
// the rightInfo format: {"add_field":["data","img"],"remove_filed":["field1"],"user_name":"qqtou-zhangsan"}
// add_field: the field which fields can be viewed
// remove_filed: the field which fields can not be viewed
// user_name: account which is registed by user chaincode
func managerFieldRight(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if err := chaincode.CheckArgsEmpty(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	rightInfo := args[0]

	mr, err := getRights(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get the right, err:%s", err.Error()))
	}

	if err := mr.managerFieldRight(stub, rightInfo); err != nil {
		return shim.Error(fmt.Sprintf("failed to manager field right, err:%s", err.Error()))
	}

	return shim.Success(nil)
}

func queryRightInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	rights, err := getRights(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get rights, err:%s", err.Error()))
	}

	if rights == nil {
		return shim.Success(nil)
	}

	if err := validateOwnerRight(stub); err == nil {
		byteRights, err := json.Marshal(rights)
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to validate owner right when json marshal right, err:%s", err.Error()))
		}
		return shim.Success(byteRights)
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get msp id by user chaincode, err:%s", err.Error()))
	}

	if !user.IsRoot && !user.IsAdmin {
		return shim.Error("unable to obtain permission information without administrator permissions")
	}

	byteRights, err := json.Marshal(rights)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to json marshal rights, err:%s", err.Error()))
	}
	return shim.Success(byteRights)
}
