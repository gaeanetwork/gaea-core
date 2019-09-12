package chaincode

import (
	"testing"

	"github.com/hyperledger/fabric/peer/common"
	"github.com/stretchr/testify/assert"
)

func Test_Package(t *testing.T) {
	ReadViperConfiguration()
	conf, err := GetConfig("system01")
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	common.InitCmd(nil, []string{})
	common.SetOrdererEnv(nil, []string{})

	conf.ChaincodePath = "chaincode/system/user"
	conf.ChaincodeName = "user"
	conf.ChaincodeVersion = "1.0"

	err = chaincodePackage(conf)
	assert.NoError(t, err)
}
