package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type RoutesConfig struct {
	AuthServicePath string `env:"AUTH_SERVICE_PATH" env-default:"http://localhost:8081"`
}

type Config struct {
	Port         int           `env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	GHTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"15s"`

	RoutesConfig
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
