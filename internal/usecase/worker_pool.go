package usecase

import (
	"context"
	"log/slog"
	"sync"

	"github.com/folivorra/ziper/internal/model"
)

type WorkerPool struct {
	ctx        context.Context
	tasks      chan *model.Task
	workersNum int
	service    *TaskService
	wg         *sync.WaitGroup
	logger     *slog.Logger
}

func NewWorkerPool(
	ctx context.Context,
	workersNum int,
	service *TaskService,
	logger *slog.Logger,
	tasks chan *model.Task,
) *WorkerPool {
	return &WorkerPool{
		ctx:        ctx,
		tasks:      tasks,
		workersNum: workersNum,
		service:    service,
		wg:         &sync.WaitGroup{},
		logger:     logger,
	}
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
	wp.logger.Info("worker pool stopped")
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
