package repository

import (
	"github.com/folivorra/ziper/internal/model"
	"sync"
)

type InMemoryTaskRepo struct {
	mu    sync.Mutex
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
