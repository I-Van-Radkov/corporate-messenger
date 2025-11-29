package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"
	v1 "github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/controller/http/v1"
)

type App struct {
	httpServer *v1.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	server := v1.NewServer(cfg)
	err := server.RegisterHandlers()
	if err != nil {
		return nil, fmt.Errorf("failed to register handlers: %w", err)
	}

	return &App{
		httpServer: server,
	}, nil
}

func (a *App) MustRun(ctx context.Context, timeout time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.httpServer.Start(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	a.httpServer.Stop(ctx)

	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		panic(err)
	}

	wg.Wait()
	log.Println("Сервис api-gateway остановлен")
}
