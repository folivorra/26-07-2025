package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/folivorra/ziper/app"
	"github.com/folivorra/ziper/internal/adapter/archiver"
	"github.com/folivorra/ziper/internal/adapter/downloader"
	"github.com/folivorra/ziper/internal/config"
	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/repository"
	"github.com/folivorra/ziper/internal/transport/rest"
	"github.com/folivorra/ziper/internal/transport/validation"
	"github.com/folivorra/ziper/internal/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

	cfg := config.NewConfig()

	a := app.NewApp(logger)
	defer a.Shutdown()

	z := archiver.NewZipArchiver(a, logger)
	d := downloader.NewHTTPDownloader(a, logger, cfg.Timeout)
	v := validation.NewHTTPValidator(cfg.Timeout)

	l := usecase.NewLockTaskManager()

	repo := repository.NewInMemoryTaskRepo()

	taskQueue := make(chan *model.Task, 3)

	ts := usecase.NewTaskService(repo, cfg.MaxTasks, cfg.MaxFilesInTask, logger, l, v, d, z, taskQueue)

	wp := usecase.NewWorkerPool(ctx, a, 3, ts, logger, taskQueue)
	wp.Start()

	srv := rest.NewServer(a, ts, logger, cfg.Port)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("failed to start http server", slog.String("error", err.Error()))
			a.Shutdown()
		}
	}()

	a.Run()
}
