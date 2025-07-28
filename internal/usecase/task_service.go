package usecase

import (
	"errors"
	"fmt"
	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/repository"
	"log/slog"
	net "net/url"
	"path"
	"sync/atomic"
)

type TaskService struct {
	repo            repository.TaskRepo
	activeTasks     atomic.Uint64
	maxTasks        uint64
	maxFilesInTasks uint64
	idCounter       atomic.Uint64
	logger          *slog.Logger
	lockManager     *LockTaskManager
}

func NewTaskService(repo repository.TaskRepo, maxTasks uint64, maxFilesInTask uint64, logger *slog.Logger) *TaskService {
	return &TaskService{
		repo:            repo,
		maxTasks:        maxTasks,
		maxFilesInTasks: maxFilesInTask,
		logger:          logger,
		lockManager:     NewLockTaskManager(),
	}
}

func (s *TaskService) CreateTask() (uint64, error) {
	if !CanAddTask(&s.activeTasks, s.maxTasks) {
		s.logger.Error("active tasks exceeds max tasks", s.maxTasks)
		return 0, fmt.Errorf("active tasks exceeds max tasks %d", s.maxTasks)
	}

	id := s.idCounter.Add(1)
	task := &model.Task{
		ID:     id,
		Status: model.TaskStatusAccepted,
		Files:  make([]*model.File, 0, 3),
	}
	s.repo.Save(task)
	s.activeTasks.Add(1)

	s.logger.Info("task created")

	return id, nil
}

func (s *TaskService) AddFileByID(id uint64, url string) error {
	lock := s.lockManager.GetLock(id)
	lock.Lock()
	defer lock.Unlock()

	status := model.FileStatusAccepted

	if _, err := net.ParseRequestURI(url); err != nil {
		s.logger.Warn("invalid url", url)
		status = model.FileStatusFailed
		//return fmt.Errorf("invalid url %s", url)
	}

	if !IsAllowedFileType(url) {
		s.logger.Warn("invalid file type", path.Ext(url))
		status = model.FileStatusInvalidType
		//return fmt.Errorf("invalid file type: %s", path.Ext(url))
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("not found task by id", id)
		return fmt.Errorf("not found task by id %d", id)
	}

	if !CanAddFileInTask(uint64(len(task.Files)), s.maxFilesInTasks) {
		s.logger.Error("task exceeds max files", s.maxFilesInTasks)
		return fmt.Errorf("task exceeds max files %d", s.maxFilesInTasks)
	}

	file := &model.File{
		Status: status,
		URL:    url,
	}

	task.Files = append(task.Files, file)
	s.logger.Info("add file to task", file)

	return nil
}
