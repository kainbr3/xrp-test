package logger

import (
	k "crypto-braza-tokens-api/utils/keys-values"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func NewLogger() {
	config := zap.NewProductionConfig()
	config.Level = getLevel()
	buildConfig, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = buildConfig
}

func getLevel() zap.AtomicLevel {
	defaultLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	level, err := k.Get("LOG_LEVEL")
	if err != nil {
		return defaultLevel
	}

	switch level {
	case "DEBUG":
		return zap.NewAtomicLevelAt(zap.DebugLevel)

	case "INFO":
		return defaultLevel

	case "WARN":
		return zap.NewAtomicLevelAt(zap.WarnLevel)

	case "ERROR":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)

	default:
		return defaultLevel
	}
}
