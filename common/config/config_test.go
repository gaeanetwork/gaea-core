package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Initialize(t *testing.T) {
	ProfileEnabled, LogLevel = true, "debug"
	assert.True(t, ProfileEnabled)
	assert.Equal(t, LogLevel, "debug")
	Initialize()
	assert.False(t, ProfileEnabled)
	assert.Equal(t, LogLevel, "info")
}
