package config

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common/glog"
	"github.com/stretchr/testify/assert"
)

func Test_Initialize(t *testing.T) {
	ProfileEnabled, glog.LogLevel = true, "info"
	assert.True(t, ProfileEnabled)
	assert.Equal(t, glog.LogLevel, "info")
	Load()
	assert.False(t, ProfileEnabled)
	assert.Equal(t, glog.LogLevel, "debug")
}
