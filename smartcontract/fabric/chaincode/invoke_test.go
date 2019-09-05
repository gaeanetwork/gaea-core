package chaincode

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/stretchr/testify/assert"
)

func Test_Invoke(t *testing.T) {
	config.Initialize()
	ReadViperConfiguration()
	conf, err := GetConfig("system01")
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	common.InitCmd(nil, []string{})
	common.SetOrdererEnv(nil, []string{})
	conf.ChaincodeInput = []string{"put", "a", "fasd"}
	_, err1 := invoke(conf)
	assert.NoError(t, err1)
}
