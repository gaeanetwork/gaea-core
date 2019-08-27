package asset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gitlab.com/jaderabbit/go-rabbit/chaincode"
	"gitlab.com/jaderabbit/go-rabbit/chaincode/system"
	"gitlab.com/jaderabbit/go-rabbit/common"
	"gitlab.com/jaderabbit/go-rabbit/common/crypto"
)

const (
	rightsKey            = "asset_right_control_stub_put_key"
	theOwnerMspIDofAsset = "asset_right_control_chaincode_instatiate_upgrade_operate_msp_id"
)

// Asset is a digital asset on chain
type Asset struct {
	Key      string            `bson:"key" json:"key"`
	MapData  map[string]*Field `bson:"map_data" json:"map_data,omitempty"`
	IsPublic bool              `bson:"is_public" json:"is_public"`
	TxID     string            `bson:"tx_id" json:"tx_id"`
	UserID   string            `bson:"user_id" json:"user_id"`
}

type Field struct {
	Data        string      `bson:"data" json:"data"`
	EncryptType crypto.Hash `bson:"encrypt_type" json:"encrypt_type"`
	IsPublic    bool        `bson:"is_public" json:"is_public"`
}

// NewAsset constructs and returns a asset instance, the Crypto value default is 0
func NewAsset() *Asset {
	return &Asset{
		MapData: make(map[string]*Field),
	}
}

func (asset *Asset) filter(stub shim.ChaincodeStubInterface) error {
	if asset.IsPublic {
		return nil
	}

	// the asset chaincode installer can view all fields
	if err := validateOwnerRight(stub); err == nil {
		return nil
	}

	user, err := getUserByChaincode(stub)
	if err != nil {
		return fmt.Errorf("failed to get user by chaincode, key:%s, txID:%s", asset.Key, stub.GetTxID())
	}

	// the asset creator can view all fields
	if user.ID == asset.UserID {
		return nil
	}

	mr, err := getRights(stub)
	if err != nil {
		return err
	}

	mapdata := make(map[string]*Field)
	if mr == nil {
		asset.MapData = mapdata
		return nil
	}

	fields, ok := mr.Fields[user.UserName]
	if !ok {
		// when no right, set the data is nil
		asset.MapData = mapdata
		return nil
	}

	if user.IsAdmin {
		if _, ok := mr.AdminUserFields[user.UserName]; ok {
			fields = append(fields, mr.AdminUserFields[user.UserName]...)
		}
	}

	mapFields := common.ConvertArrayToMap(fields)
	for fieldName, fieldInfo := range asset.MapData {
		if fieldInfo.IsPublic {
			mapdata[fieldName] = fieldInfo
			continue
		}

		if _, ok := mapFields[fieldName]; ok {
			mapdata[fieldName] = fieldInfo
		}
	}

	asset.MapData = mapdata
	return nil
}

func (f *Field) Equal(compareField *Field) (bool, error) {
	if f.EncryptType == crypto.SM2 {
		return false, errors.New("Sm2 encryption algorithm, please obtain cipher text, after decryption and original text comparison to verify")
	}

	if f.Data != compareField.Data {
		return false, nil
	}

	if f.EncryptType != compareField.EncryptType {
		return false, nil
	}

	return true, nil
}

func validateOwnerRight(stub shim.ChaincodeStubInterface) error {
	mspID, err := chaincode.GetMSPID(stub)
	if err != nil {
		return fmt.Errorf("failed to get msp id, err:%s", err.Error())
	}

	byteOwner, err := stub.GetState(theOwnerMspIDofAsset)
	if err != nil {
		return fmt.Errorf("failed to get the state of theOwnerMspIDofAsset, err:%s", err.Error())
	}

	if mspID != string(byteOwner) {
		return fmt.Errorf("failed to validate owner right")
	}
	return nil
}

func getUserByChaincode(stub shim.ChaincodeStubInterface) (*system.User, error) {
	res := stub.InvokeChaincode("user", [][]byte{[]byte("getuser")}, "")
	if res.Status != shim.OK {
		return nil, errors.New(res.Message)
	}

	user := &system.User{}
	if err := json.Unmarshal(res.Payload, user); err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByUserName(stub shim.ChaincodeStubInterface, userName string) (*system.User, error) {
	resID := stub.InvokeChaincode("user", [][]byte{[]byte("query"), []byte("user"), []byte(userName)}, "")
	if resID.Status != shim.OK {
		return nil, errors.New(resID.Message)
	}

	res := stub.InvokeChaincode("user", [][]byte{[]byte("query"), []byte(""), resID.Payload}, "")
	user := &system.User{}
	if err := json.Unmarshal(res.Payload, user); err != nil {
		return nil, err
	}
	return user, nil
}
