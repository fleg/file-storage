package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"dev" validate:"oneof=dev prod"`
	Server struct {
		Port uint `env:"PORT" envDefault:"9000" validate:"gt=0,lt=65535"`
	}
	Logger struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug" validate:"oneof=debug error info"`
	}
	PostgreSQL struct {
		Url string `env:"POSTGRESQL_URL"`
	}
	Storage struct {
		Path        string `env:"STORAGE_PATH" envDefault:"/tmp"`
		MaxFileSize uint   `env:"STORAGE_MAX_FILE_SIZE" envDefault:"134217728"`
	}
}

func Load(c *Config) error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		return fmt.Errorf("validation of config failed: %w", err)
	}

	return nil
}
