package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	shutdownCh chan os.Signal
	cleanup    []func(context.Context)
	logger     *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		shutdownCh: make(chan os.Signal, 1),
		cleanup:    []func(context.Context){},
		logger:     logger,
	}
}

func (a *App) Run() {
	signal.Notify(a.shutdownCh, syscall.SIGTERM, os.Interrupt)

	a.logger.Info("app started")

	<-a.shutdownCh
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.logger.Info("app shutting down")

	for i := len(a.cleanup) - 1; i >= 0; i-- {
		a.cleanup[i](ctx)
	}

	a.logger.Info("app shutdown complete")
}

func (a *App) RegisterCleanup(f func(context.Context)) {
	a.cleanup = append(a.cleanup, f)
}
