package glog

import (
	"time"

	"github.com/gaeanetwork/gaea-core/common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MustGetLogger get a gaea logger
func MustGetLogger() *zap.Logger {
	logger, err := newZapConfig(NameToLevel(config.LogLevel)).Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func newZapConfig(level zapcore.Level) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console", // json cannot set color
		EncoderConfig:    newZapCoreEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newZapCoreEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
