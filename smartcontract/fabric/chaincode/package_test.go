package chaincode

import (
	"testing"

	"github.com/hyperledger/fabric/peer/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
)

func Test_Package(t *testing.T) {
	config.Initialize()
	ReadViperConfiguration()
	conf, err := GetConfig("system01")
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	common.InitCmd(nil, []string{})
	common.SetOrdererEnv(nil, []string{})
	config.Initialize()

	conf.ChaincodePath = "gitlab.com/jaderabbit/go-rabbit/chaincode/system/user"
	conf.ChaincodeName = "user"
	conf.ChaincodeVersion = "1.0"

	err = chaincodePackage(conf)
	assert.NoError(t, err)
}
