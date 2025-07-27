package repository

import (
	"errors"
	"github.com/folivorra/ziper/internal/model"
	"log/slog"
	"sync"
	"sync/atomic"
)

type InMemoryTaskRepo struct {
	mu          sync.Mutex
	tasks       map[uint64]*model.Task
	activeTasks int
	maxTasks    int
	counter     atomic.Uint64
	logger      *slog.Logger
}

func NewTaskManager(maxTasks int, logger slog.Logger) *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks:    make(map[uint64]*model.Task),
		maxTasks: maxTasks,
		counter:  atomic.Uint64{},
		logger:   &logger,
	}
}

func (t *InMemoryTaskRepo) CreateTask() (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.activeTasks >= t.maxTasks {
		t.logger.Error("too many active tasks")
		return 0, errors.New("too many active tasks")
	}

	id := t.counter.Add(1)
	task := &model.Task{
		ID:     id,
		Status: model.TaskStatusAccepted,
		Files:  make([]model.File, 0, 3),
	}
	t.tasks[id] = task
	t.activeTasks++

	t.logger.Info("task created")

	return id, nil
}
