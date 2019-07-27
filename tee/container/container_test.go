package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetContainer(t *testing.T) {
	c, err := GetContainer(Dev)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c1, err := GetContainer(Docker)
	assert.NoError(t, err)
	assert.NotNil(t, c1)

	_, err = GetContainer(SGX)
	assert.Error(t, err)

	_, err = GetContainer(Type(954))
	assert.Error(t, err)
}
