package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cleanup    []func(context.Context)
	shutdownCh chan os.Signal
	logger     *slog.Logger
}

func NewApp(ctx context.Context, logger slog.Logger) *App {
	newCtx, cancel := context.WithCancel(ctx)
	return &App{
		ctx:        newCtx,
		shutdownCh: make(chan os.Signal),
		cancel:     cancel,
		logger:     &logger,
	}
}

func (a *App) Run() {
	signal.Notify(a.shutdownCh, syscall.SIGTERM, syscall.SIGINT)

	a.logger.Info("app started")

	<-a.shutdownCh
}

func (a *App) Shutdown() {
	a.logger.Info("app shutting down")
	for i := len(a.cleanup) - 1; i >= 0; i-- {
		a.cleanup[i](a.ctx)
	}
	a.cancel()
	a.logger.Info("app shut down")
}

func (a *App) RegisterCleanup(f func(context.Context)) {
	a.cleanup = append(a.cleanup, f)
}
