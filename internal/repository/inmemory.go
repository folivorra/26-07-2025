package repository

import (
	"fmt"
	"github.com/folivorra/ziper/internal/model"
	"sync"
)

type InMemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[uint64]*model.Task
}

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
