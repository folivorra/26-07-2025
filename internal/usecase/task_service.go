package usecase

import (
	"fmt"
	"log/slog"
	net "net/url"
	"path"
	"sync"
	"sync/atomic"

	"github.com/folivorra/ziper/internal/adapter/archiver"
	"github.com/folivorra/ziper/internal/adapter/downloader"
	"github.com/folivorra/ziper/internal/config"
	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/repository"
	"github.com/folivorra/ziper/internal/transport/validation"
)

type TaskService struct {
	repo        repository.TaskRepo
	activeTasks atomic.Uint64
	cfg         config.Config
	idCounter   atomic.Uint64
	lockManager *LockTaskManager
	validr      validation.FileValidator
	dowloadr    downloader.Downloader
	archiver    archiver.Archiver
	logger      *slog.Logger
	taskQueue   chan *model.Task
}

func NewTaskService(
	repo repository.TaskRepo,
	cfg config.Config,
	logger *slog.Logger,
	locker *LockTaskManager,
	validr validation.FileValidator,
	dowloadr downloader.Downloader,
	archiver archiver.Archiver,
	taskQueue chan *model.Task,
) *TaskService {
	return &TaskService{
		repo:        repo,
		cfg:         cfg,
		lockManager: locker,
		validr:      validr,
		dowloadr:    dowloadr,
		archiver:    archiver,
		logger:      logger,
		taskQueue:   taskQueue,
	}
}

func (s *TaskService) CreateTask() (uint64, error) {
	s.logger.Info("creating new task")

	for {
		current := s.activeTasks.Load()
		if current >= s.cfg.MaxTasks {
			s.logger.Error("active tasks exceeds max tasks",
				slog.Uint64("max tasks", s.cfg.MaxTasks),
				slog.Uint64("active tasks", current),
			)
			return 0, fmt.Errorf("active tasks exceeds max tasks %d", s.cfg.MaxTasks)
		}
		if s.activeTasks.CompareAndSwap(current, current+1) {
			break
		}
	}

	id := s.idCounter.Add(1)
	task := &model.Task{
		ID:          id,
		Status:      model.TaskStatusAccepted,
		Files:       make([]*model.File, 0, s.cfg.MaxFilesInTask),
		ArchiveURL:  fmt.Sprintf("http://localhost:%s/%s/task-%d.zip", s.cfg.Port, s.cfg.ArchDir, id),
		ArchivePath: fmt.Sprintf("%s/task-%d.zip", s.cfg.ArchDir, id),
	}
	s.repo.Save(task)

	s.logger.Info("task created")

	return id, nil
}

func (s *TaskService) AddFileByID(id uint64, url string) (model.FileStatus, error) {
	lock := s.lockManager.GetLock(id)
	lock.Lock()
	defer lock.Unlock()

	s.logger.Info("adding file to task",
		slog.Uint64("id", id),
		slog.String("url", url),
	)

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("error getting task by id",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return model.FileStatusFailed, fmt.Errorf("not found task by id %d", id)
	}

	if !CanAddFileInTask(uint64(len(task.Files)), s.cfg.MaxFilesInTask) {
		s.logger.Error("task exceeds max files",
			slog.Uint64("maxFilesInTask", s.cfg.MaxFilesInTask),
			slog.Uint64("currentFiles", uint64(len(task.Files))),
		)
		return model.FileStatusFailed, fmt.Errorf("task exceeds max files %d", s.cfg.MaxFilesInTask)
	}

	status := model.FileStatusAccepted
	var returningErr error

	if _, err := net.ParseRequestURI(url); err != nil {
		s.logger.Warn("invalid url",
			slog.String("url", url),
			slog.String("error", err.Error()),
		)
		status = model.FileStatusInvalidURL
		returningErr = fmt.Errorf("invalid url %s", url)
	} else if !IsAllowedFileType(url) {
		s.logger.Warn("not supported file type",
			slog.String("file type", path.Ext(url)),
		)
		status = model.FileStatusNotSupportedType
		returningErr = fmt.Errorf("not supported file type %s", path.Ext(url))
	} else if !s.validr.IsReachable(url) {
		s.logger.Warn("file not reachable",
			slog.String("url", url),
		)
		status = model.FileStatusNotReachable
		returningErr = fmt.Errorf("file not reachable %s", url)
	}

	file := &model.File{
		Status: status,
		URL:    url,
	}

	task.Files = append(task.Files, file)

	if len(task.Files) == int(s.cfg.MaxFilesInTask) {
		s.taskQueue <- task
		s.logger.Info("task goes to queue")
	}

	s.logger.Info("added file to task",
		slog.Uint64("task_id", task.ID),
		slog.String("file_status", string(file.Status)),
		slog.String("file_url", file.URL),
	)

	return status, returningErr
}

func (s *TaskService) GetTaskStatusAndArchiveURL(id uint64) (model.TaskStatus, string, error) {
	s.logger.Info("getting task status",
		slog.Uint64("id", id),
	)

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("error getting task by id",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return "", "", fmt.Errorf("not found task by id %d", id)
	}

	status := task.Status
	archURL := ""
	if len(task.Files) == int(s.cfg.MaxFilesInTask) || task.Status == model.TaskStatusCompleted {
		archURL = task.ArchiveURL
		s.logger.Info("got archive url",
			slog.Uint64("task_id", task.ID),
		)
	}

	s.logger.Info("got task status",
		slog.Uint64("task_id", task.ID),
	)

	return status, archURL, nil
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

	dirPath := fmt.Sprintf("downloads/task-%d", task.ID)

	sem := NewSemaphore(int(s.cfg.MaxFilesInTask))
	var wg sync.WaitGroup

	for _, file := range task.Files {
		sem.Acquire()
		wg.Add(1)
		go func(file *model.File) {
			defer wg.Done()
			defer sem.Release()
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("panic during file processing",
						slog.Uint64("task_id", task.ID),
						slog.String("file_url", file.URL),
						slog.Any("error", r),
					)
				}
				file.Status = model.FileStatusFailed
			}()

			s.logger.Info("downloading file",
				slog.Uint64("task_id", task.ID),
				slog.String("file_url", file.URL),
			)

			err := s.dowloadr.DownloadFile(file.URL, task.ID)
			if err != nil {
				s.logger.Error("error downloading file",
					slog.Uint64("task_id", task.ID),
					slog.String("file_url", file.URL),
					slog.String("error", err.Error()),
				)
				file.Status = model.FileStatusFailed
			} else {
				file.Status = model.FileStatusCompleted
				s.logger.Info("file downloading successfully",
					slog.Uint64("task_id", task.ID),
					slog.String("file_url", file.URL),
				)
			}
		}(file)
	}

	wg.Wait()

	if err := s.archiver.ArchiveDirectory(dirPath); err != nil {
		s.logger.Error("error adding file to archive",
			slog.Uint64("task_id", task.ID),
			slog.String("dir_path", dirPath),
			slog.String("error", err.Error()),
		)
	}

	task.Status = model.TaskStatusCompleted
	s.activeTasks.Add(^uint64(0))

	s.logger.Info("task processing completed",
		slog.Uint64("id", task.ID),
		slog.String("status", string(task.Status)),
	)

	return nil
}
