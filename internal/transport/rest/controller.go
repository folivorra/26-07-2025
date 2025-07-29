package rest

import (
	"github.com/folivorra/ziper/internal/usecase"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Controller struct {
	taskService *usecase.TaskService
	logger      *slog.Logger
}

func NewController(taskService *usecase.TaskService, logger *slog.Logger) *Controller {
	return &Controller{
		taskService: taskService,
		logger:      logger,
	}
}

func (c *Controller) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) AddFileByIDHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) GetTaskStatusAndArchivePathHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) DownloadArchiveHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", c.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", c.GetTaskStatusAndArchivePathHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}/add", c.AddFileByIDHandler).Methods("POST")
	r.HandleFunc("/archives/{filename}", c.DownloadArchiveHandler).Methods("GET")
}
