package logger

import (
	"fmt"

	"example.com/sample/commonlib/config"
	"example.com/sample/commonlib/config/sharedoptions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger for a go microservice
var Log *zap.Logger

// logLevel is an atomic container which allows the real-time updating of the log level for Log
var logLevel *zap.AtomicLevel

// LevelFromString converts a log level string to a zapcore.Level, returning an error for unrecognized levels
func LevelFromString(levelStr string) (zapcore.Level, error) {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return 0, fmt.Errorf("not a valid log level: %v", levelStr)
	}
}

// InitLogger initializes the global logger, Log. It accepts parameters for whether
// the application is running in development mode and the initial log level
func InitLogger(level zapcore.Level, isProduction bool) error {
	var configuration zap.Config
	if isProduction {
		configuration = zap.NewProductionConfig()
	} else {
		configuration = zap.NewDevelopmentConfig()
	}
	logLevel = &configuration.Level
	logLevel.SetLevel(level)

	var err error
	Log, err = configuration.Build()
	return err
}

// InitLoggerFromConfig initializes the global logger, Log, via shared options in a config.Registry. Notably,
// it requires the registry to have the options sharedoptions.LogLevel and sharedoptions.IsInProduction registered.
func InitLoggerFromConfig(registry config.Registry) error {
	logLevelStr, levelPresent := registry.Get(sharedoptions.LogLevel)
	if !levelPresent {
		logLevelStr = "info"
	}
	parsedLevel, parseErr := LevelFromString(logLevelStr)
	if parseErr != nil {
		return fmt.Errorf("invalid log level bypassed validation: %v cause: %w", logLevelStr, parseErr)
	}
	isProductionStr := registry.GetRequired(sharedoptions.IsInProduction)

	if loggerSetupErr := InitLogger(parsedLevel, isProductionStr == "true"); loggerSetupErr != nil {
		return fmt.Errorf("logger setup failed: %w", loggerSetupErr)
	}

	return nil
}

// AdjustLevel updates the log level for the global logger, Log
func AdjustLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}
