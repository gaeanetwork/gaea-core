package chaincode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConfig(t *testing.T) {
	conf, err := GetConfig("")
	assert.NotNil(t, err)
	assert.Nil(t, conf)

	conf, err = GetConfig("debt111")
	assert.NotNil(t, err)
	assert.Nil(t, conf)

	conf, err = GetConfig("debt")
	assert.Nil(t, err)
	assert.NotNil(t, conf)
}
