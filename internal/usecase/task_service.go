package usecase

import (
	"fmt"
	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/repository"
	"log/slog"
	"sync/atomic"
)

type TaskService struct {
	repo        repository.TaskRepo
	activeTasks atomic.Uint64
	maxTasks    uint64
	idCounter   atomic.Uint64
	logger      *slog.Logger
}

func NewTaskService(repo repository.TaskRepo, maxTasks uint64, logger *slog.Logger) *TaskService {
	return &TaskService{
		repo:     repo,
		maxTasks: maxTasks,
		logger:   logger,
	}
}

func (s *TaskService) CreateTask() (uint64, error) {
	if s.activeTasks.Load() >= s.maxTasks {
		return 0, fmt.Errorf("active tasks exceeds max tasks %d", s.maxTasks)
	}

	id := s.idCounter.Add(1)
	task := &model.Task{
		ID:     id,
		Status: model.TaskStatusAccepted,
		Files:  make([]model.File, 0, 3),
	}
	s.repo.Save(task)
	s.activeTasks.Add(1)
	
	return id, nil
}
