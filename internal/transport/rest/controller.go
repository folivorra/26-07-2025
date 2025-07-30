package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"github.com/folivorra/ziper/internal/model"
	"github.com/folivorra/ziper/internal/usecase"
	"github.com/gorilla/mux"
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

func (c *Controller) CreateTaskHandler(w http.ResponseWriter, _ *http.Request) {
	id, err := c.taskService.CreateTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		ID uint64 `json:"id"`
	}{
		ID: id,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) AddFileByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request := struct {
		URL string `json:"url"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.taskService.AddFileByID(id, request.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetTaskStatusAndArchivePathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status, url, err := c.taskService.GetTaskStatusAndArchiveURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status model.TaskStatus `json:"status"`
		URL    string           `json:"path,omitempty"`
	}{
		Status: status,
		URL:    url,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) DownloadArchiveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	idStr := filename[len("task-") : len(filename)-len(".zip")]
	id, err := strconv.ParseUint(idStr, 10, 64)

	_, archiveURL, err := c.taskService.GetTaskStatusAndArchiveURL(id)
	if err != nil {
		http.Error(w, "failed to get archive path", http.StatusInternalServerError)
		return
	}

	if archiveURL == "" {
		http.Error(w, "archive not ready", http.StatusNotFound)
		return
	}

	filePath := path.Join("archives", filename)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"task-%d.zip\"", id))
	http.ServeFile(w, r, filePath)
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", c.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", c.GetTaskStatusAndArchivePathHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}/add", c.AddFileByIDHandler).Methods("POST")
	r.HandleFunc("/archives/{filename:.+}", c.DownloadArchiveHandler).Methods("GET")
}
