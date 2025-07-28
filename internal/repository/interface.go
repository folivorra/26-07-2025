package repository

import "github.com/folivorra/ziper/internal/model"

type TaskRepo interface {
	Save(task *model.Task)
	GetByID(id uint64) (*model.Task, error)
	GetTaskStatus(id uint64) (model.TaskStatus, error)
}
