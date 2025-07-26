package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cleanup    []func(context.Context)
	shutdownCh chan os.Signal
}

func NewApp(ctx context.Context) *App {
	newCtx, cancel := context.WithCancel(ctx)
	return &App{
		ctx:        newCtx,
		shutdownCh: make(chan os.Signal),
		cancel:     cancel,
	}
}

func (a *App) Run() {
	signal.Notify(a.shutdownCh, syscall.SIGTERM, syscall.SIGINT)

	<-a.shutdownCh
}

func (a *App) Shutdown() {
	for i := len(a.cleanup) - 1; i >= 0; i-- {
		a.cleanup[i](a.ctx)
	}
	a.cancel()
}

func (a *App) RegisterCleanup(f func(context.Context)) {
	a.cleanup = append(a.cleanup, f)
}
