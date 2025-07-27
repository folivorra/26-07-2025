package usecase

type TaskService struct {
	repo TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (t *TaskService) CreateTask() (uint64, error) {
	return t.repo.CreateTask()
}
