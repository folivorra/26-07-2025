package usecase

type TaskRepo interface {
	CreateTask() (uint64, error)
}
