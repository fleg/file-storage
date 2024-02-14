package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger = zap.SugaredLogger

	// TODO: env as enum or smth?
	Options struct {
		Env   string
		Level string
	}
)

func getZapLogLevel(level string) (zapcore.Level, error) {
	var zapLogLevel zapcore.Level

	err := zapLogLevel.UnmarshalText([]byte(level))
	return zapLogLevel, err
}

func New(opt Options) (*Logger, error) {
	level, err := getZapLogLevel(opt.Level)
	if err != nil {
		return nil, err
	}

	var config zap.Config
	switch opt.Env {

	case "prod":
		config = zap.NewProductionConfig()

	case "dev":
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(level)

	default:
		return nil, fmt.Errorf("logger unsupported env: %q", opt.Env)
	}

	config.InitialFields = map[string]interface{}{
		"env": opt.Env,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
