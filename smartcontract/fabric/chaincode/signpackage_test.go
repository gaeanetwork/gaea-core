package chaincode

import (
	"testing"

	"github.com/hyperledger/fabric/peer/common"
	"github.com/stretchr/testify/assert"
)

func Test_SignPackage(t *testing.T) {
	ReadViperConfiguration()
	conf, err := GetConfig("system01")
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	common.InitCmd(nil, []string{})
	common.SetOrdererEnv(nil, []string{})

	conf.ChaincodeName = "user"

	err = signpackage(conf)
	assert.NoError(t, err)
}
