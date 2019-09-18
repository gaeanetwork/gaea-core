package config

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common/glog"
	"github.com/stretchr/testify/assert"
)

func Test_Initialize(t *testing.T) {
	ProfileEnabled, glog.LogLevel = true, "debug"
	assert.True(t, ProfileEnabled)
	assert.Equal(t, glog.LogLevel, "debug")
	initialize()
	assert.False(t, ProfileEnabled)
	assert.Equal(t, glog.LogLevel, "info")
}
