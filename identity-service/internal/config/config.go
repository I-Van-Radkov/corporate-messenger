package config

import (
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/clients/directory"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/usecase"
	postgres "github.com/I-Van-Radkov/corporate-messenger/identity-service/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         int           `env:"PORT" env-default:"8081"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	GHTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"15s"`

	postgres.PostgresConfig

	directory.DirectoryServiceConfig

	usecase.AuthConfig
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
