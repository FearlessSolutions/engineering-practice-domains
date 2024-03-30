package loglevel

import (
	"example.com/sample/commonlib/logger"
	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -destination ./adjust_log_leveL_mocks.go -package loglevel . Core

// Core contains logic for talking to the global logger.
type Core interface {
	// SetLogLevel adjusts the global logger to the requested log level
	SetLogLevel(level zapcore.Level)
}

type CoreLogic struct{}

// SetLogLevel implements Core for CoreLogic
func (CoreLogic) SetLogLevel(level zapcore.Level) {
	logger.AdjustLevel(level)
}
