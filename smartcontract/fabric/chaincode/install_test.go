package chaincode

import (
	"io/ioutil"
	"testing"

	"github.com/hyperledger/fabric/peer/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jaderabbit/go-rabbit/core/config"
)

func Test_Install(t *testing.T) {
	config.Initialize()
	ReadViperConfiguration()
	conf, err := GetConfig("system01")
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	common.InitCmd(nil, []string{})
	common.SetOrdererEnv(nil, []string{})

	ccpackfile, err := ioutil.ReadFile("/home/wff/go/src/gitlab.com/jaderabbit/go-rabbit/project/echain/signeddepartmentpack.out")
	assert.NoError(t, err)

	err = install(conf, ccpackfile)
	assert.NoError(t, err)
}
