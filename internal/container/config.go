package container

import (
	"fmt"

	"file-storage/internal/config"
	"go.uber.org/fx"
)

func ProvideConfig() (*config.Config, error) {
	c := &config.Config{}

	if err := config.Load(c); err != nil {
		return nil, fmt.Errorf("load config failed: %w", err)
	}

	return c, nil
}

var ConfigModule = fx.Module("config", fx.Options(
	fx.Provide(ProvideConfig),
))
