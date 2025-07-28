package usecase

import (
	"fmt"
	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/repository"
	"log/slog"
	net "net/url"
	"path"
	"sync/atomic"
)

type TaskService struct {
	repo           repository.TaskRepo
	activeTasks    atomic.Uint64
	maxTasks       uint64
	maxFilesInTask uint64
	idCounter      atomic.Uint64
	logger         *slog.Logger
	lockManager    *LockTaskManager
}

func NewTaskService(repo repository.TaskRepo, maxTasks uint64, maxFilesInTask uint64, logger *slog.Logger) *TaskService {
	return &TaskService{
		repo:           repo,
		maxTasks:       maxTasks,
		maxFilesInTask: maxFilesInTask,
		logger:         logger,
		lockManager:    NewLockTaskManager(),
	}
}

func (s *TaskService) CreateTask() (uint64, error) {
	s.logger.Info("creating new task")

	if !CanAddTask(&s.activeTasks, s.maxTasks) {
		s.logger.Error("active tasks exceeds max tasks",
			slog.Uint64("max tasks", s.maxTasks),
			slog.Uint64("active tasks", s.activeTasks.Load()),
		)
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

	s.logger.Info("adding file to task",
		slog.Uint64("id", id),
		slog.String("url", url),
	)

	status := model.FileStatusAccepted

	if _, err := net.ParseRequestURI(url); err != nil {
		s.logger.Warn("invalid url",
			slog.String("url", url),
			slog.String("error", err.Error()),
		)
		status = model.FileStatusFailed
		//return fmt.Errorf("invalid url %s", url)
	}

	if !IsAllowedFileType(url) {
		s.logger.Warn("invalid file type",
			slog.String("file type", path.Ext(url)),
		)
		status = model.FileStatusInvalidType
		//return fmt.Errorf("invalid file type: %s", path.Ext(url))
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("error getting task by id",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("not found task by id %d", id)
	}

	if !CanAddFileInTask(uint64(len(task.Files)), s.maxFilesInTask) {
		s.logger.Error("task exceeds max files",
			slog.Uint64("maxFilesInTask", s.maxFilesInTask),
			slog.Uint64("currentFiles", uint64(len(task.Files))),
		)
		return fmt.Errorf("task exceeds max files %d", s.maxFilesInTask)
	}

	file := &model.File{
		Status: status,
		URL:    url,
	}

	task.Files = append(task.Files, file)
	s.logger.Info("added file to task", file)

	return nil
}

func (s *TaskService) GetTaskStatusAndZipPath(id uint64) (model.TaskStatus, string, error) {
	s.logger.Info("getting task status", slog.Uint64("id", id))

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("error getting task by id",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return "", "", fmt.Errorf("not found task by id %d", id)
	}

	status := task.Status
	//if len(task.Files) == 3 {
	//	zipPath := GetZipPath()
	//}

	return status, "", nil
}

func (s *TaskService) ProcessTask(task *model.Task) error {
	lock := s.lockManager.GetLock(task.ID)
	lock.Lock()
	defer lock.Unlock()

	s.logger.Info("processing task",
		slog.Uint64("id", task.ID),
	)

	if task.Status != model.TaskStatusAccepted {
		s.logger.Warn("task already processed",
			slog.Uint64("id", task.ID),
			slog.String("status", string(task.Status)),
		)
		return fmt.Errorf("task already processed with status %s", task.Status)
	}

	task.Status = model.TaskStatusInProgress

	sem := NewSemaphore(3)
	for _, file := range task.Files {
		sem.Acquire()
		go func(file *model.File) {
			defer sem.Release()

			s.logger.Info("processing file",
				slog.Uint64("task_id", task.ID),
				slog.String("file_url", file.URL),
			)

			//.........

			file.Status = model.FileStatusCompleted
			s.logger.Info("file processed successfully",
				slog.Uint64("task_id", task.ID),
				slog.String("file_url", file.URL),
			)
		}(file)
	}

	task.Status = model.TaskStatusCompleted

	s.logger.Info("task processing completed",
		slog.Uint64("id", task.ID),
		slog.String("status", string(task.Status)),
	)

	s.activeTasks.Add(^uint64(0))

	return nil
}
