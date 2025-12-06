package config

import (
	"fmt"
	"time"

	postgres "github.com/I-Van-Radkov/corporate-messenger/chat-service/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         int           `env:"PORT"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT"`

	GHTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT"`

	postgres.PostgresConfig
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
