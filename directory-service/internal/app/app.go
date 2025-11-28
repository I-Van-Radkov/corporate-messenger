package app

import (
	"fmt"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/config"
	v1 "github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/controller/v1"
	postgres "github.com/I-Van-Radkov/corporate-messenger/directory-service/pkg/db"
)

type App struct {
	httpServer *v1.Server
	postgresDB *postgres.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	server := v1.NewServer(cfg.Port, cfg.ReadTimeout, cfg.WriteTimeout, db.Pool)

	return &App{
		httpServer: server,
		postgresDB: db,
	}, nil
}
