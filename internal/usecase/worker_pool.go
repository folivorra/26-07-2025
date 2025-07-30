package usecase

import (
	"context"
	"log/slog"
	"sync"

	"github.com/folivorra/ziper/app"
	"github.com/folivorra/ziper/internal/model"
)

type WorkerPool struct {
	ctx        context.Context
	app        *app.App
	tasks      chan *model.Task
	workersNum int
	service    *TaskService
	wg         *sync.WaitGroup
	logger     *slog.Logger
}

func NewWorkerPool(
	ctx context.Context,
	app *app.App,
	workersNum int,
	service *TaskService,
	logger *slog.Logger,
	tasks chan *model.Task,
) *WorkerPool {
	wp := &WorkerPool{
		ctx:        ctx,
		app:        app,
		tasks:      tasks,
		workersNum: workersNum,
		service:    service,
		wg:         &sync.WaitGroup{},
		logger:     logger,
	}

	wp.app.RegisterCleanup(func(ctx context.Context) {
		wp.Stop()
		wp.logger.Info("worker pool shutdown complete")
	})

	return wp
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workersNum; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for {
				select {
				case <-wp.ctx.Done():
					return
				case task, ok := <-wp.tasks:
					if !ok {
						return
					}

					func(task *model.Task) {
						defer func() {
							if r := recover(); r != nil {
								wp.logger.Error("worker panicked",
									slog.Int("worker_id", workerID),
									slog.Any("error", r),
								)
							}
						}()

						wp.logger.Info("worker started processing task",
							slog.Int("worker_id", workerID),
							slog.Uint64("task_id", task.ID),
						)
						if err := wp.service.ProcessTask(task); err != nil {
							wp.logger.Error("error processing task",
								slog.Int("worker_id", workerID),
								slog.Uint64("task_id", task.ID),
								slog.String("error", err.Error()),
							)
						}
					}(task)
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.tasks)
	wp.wg.Wait()
}

//func (wp *WorkerPool) AddTask(task *model.Task) {
//	select {
//	case wp.tasks <- task:
//		wp.logger.Info("task added to worker pool",
//			slog.Uint64("task_id", task.ID),
//		)
//	case <-wp.ctx.Done():
//		wp.logger.Warn("context done, cannot add task to worker pool",
//			slog.Uint64("task_id", task.ID),
//		)
//	}
//}
