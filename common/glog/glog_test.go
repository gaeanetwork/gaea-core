package glog

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/stretchr/testify/assert"
)

func Test_Print(t *testing.T) {
	logger4 := MustGetLogger()
	logger4.Info("???,what's this?")
	config.LogLevel = "debug"
	logger5 := MustGetLogger()
	logger5.Debug("???,what's this?")
	logger5.Warn("???,what's this?")
	logger5.Error("???,what's this?")
	logger5.DPanic("???,what's this?")
	assert.Panics(t, func() {
		logger5.Panic("???,what's this?")
	})
	// logger5.Fatal("???,what's this?")
}
