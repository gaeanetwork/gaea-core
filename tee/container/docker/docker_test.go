package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	c, err := Create()
	assert.NoError(t, err)
	defer c.Destroy()
	assert.NotNil(t, c.client)
	assert.NotNil(t, c.id)

	// Repeat create
	c, err1 := Create()
	assert.NoError(t, err1)

	cmd := fmt.Sprintf("echo 'hello world'")
	container, err2 := c.startFunc(cmd)
	assert.NoError(t, err2)

	container1, err3 := c.client.InspectContainer(container.ID)
	assert.NoError(t, err3)
	assert.Equal(t, container1.Args[1], cmd)
}
