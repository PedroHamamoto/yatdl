package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"yatdl/internal/http/middleware"
)

type Handler struct {
	service *Service
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type taskResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	var request createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("failed to decode task request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	input := CreateTaskInput{
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
	}

	task, err := h.service.CreateTask(r.Context(), input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := taskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	taskIDParam := r.PathValue("id")
	taskID, err := strconv.ParseUint(taskIDParam, 10, 64)
	if err != nil {
		log.Printf("failed to parse task id: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}
	var request updateTaskRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("failed to decode task request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	err = h.service.UpdateTask(r.Context(), UpdateTaskInput{
		ID:          taskID,
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
		Completed:   request.Completed,
	})

	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case errors.Is(err, ErrCannotUpdateTaskFromAnotherUser):
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
