package repository

import (
	"fmt"
	"sync"

	"github.com/folivorra/ziper/internal/model"
)

type InMemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[uint64]*model.Task
}

var _ TaskRepo = (*InMemoryTaskRepo)(nil)

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks: make(map[uint64]*model.Task),
	}
}

func (t *InMemoryTaskRepo) Save(task *model.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasks[task.ID] = task
}

func (t *InMemoryTaskRepo) GetByID(id uint64) (*model.Task, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	task, ok := t.tasks[id]

	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}
