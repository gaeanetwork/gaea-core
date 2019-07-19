package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	container := New()
	assert.NotNil(t, container)
}
