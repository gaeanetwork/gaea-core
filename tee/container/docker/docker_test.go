package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func Test_Create(t *testing.T) {
	c := New()
	err := c.Create()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.address)
}
