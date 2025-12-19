package config

import (
	"fmt"
	"os"
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
	var cfg Config

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}

	if err := cleanenv.ReadConfig(envPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", envPath, err)
	}

	return &cfg, nil
}
