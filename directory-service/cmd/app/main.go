package main

import (
	"context"
	"fmt"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/app"
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/config"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to creating the app structure: %w", err))
	}

	ctx := context.Background()
	app.MustRun(ctx, cfg.GHTimeout)
}
