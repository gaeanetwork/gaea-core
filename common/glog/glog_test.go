package glog

import (
	"errors"
	"testing"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func Test_WithFields(t *testing.T) {
	logger := MustGetLogger()
	logger.Info("???,what's this?", zap.String("hello", "world"), zap.String("hello", "world"))
	logger.With(zap.String("hello", "world"), zap.String("failure", "oh no"),
		zap.Stack("stack"), zap.Int("count", 42)).Info("???,what's this?")
	// , zap.Object("user", user.User{Username: "alice"})

	logger.Info("???,what's this?")
	logger.Named("glog").Info("???,what's this?")
	MustGetLoggerWithNamed("glog").Info("???,what's this?")
	MustGetLoggerWithNamedAndModule("glog", "glog").Info("???,what's this?")
	logger.Sugar().Infof("???,what's this? %s", "format you can use this")
	logger.Sugar().With("hello", "world", "failure", errors.New("oh no"),
		zap.Stack("stack"), "count", 42).Info("???,what's this?")
	// , "user", User{Name: "alice"}
}
