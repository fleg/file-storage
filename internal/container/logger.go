package container

import (
	"context"
	"fmt"

	"go-files/internal/config"
	"go-files/internal/logger"
	"go.uber.org/fx"
)

func ProvideLogger(c *config.Config) (*logger.Logger, error) {
	l, err := logger.New(
		logger.Options{
			Env:   c.AppEnv,
			Level: c.Logger.Level,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("logger couldn't be inialized: %w", err)
	}

	return l, nil
}

func RegisterLoggerHooks(lc fx.Lifecycle, logger *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sync()
			return nil
		},
	})
}

var LoggerModule = fx.Options(
	fx.Provide(ProvideLogger),
	fx.Invoke(RegisterLoggerHooks),
)
