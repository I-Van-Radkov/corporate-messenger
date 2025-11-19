package app

import (
	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"
	v1 "github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/controller/v1"
)

type App struct {
	httpServer *v1.Server
}

func NewApp(cfg *config.Config) *App {
	server := v1.NewServer(cfg.Port, cfg.ReadTimeout, cfg.WriteTimeout)
	return &App{}
}
