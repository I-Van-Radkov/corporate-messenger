package main

import (
	"context"
	"fmt"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/app"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/config"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	fmt.Println(cfg)

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to creating the app structure: %w", err))
	}

	ctx := context.Background()
	app.MustRun(ctx, cfg.GHTimeout)
}
