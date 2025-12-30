package task

import (
	"encoding/json"
	"log"
	"net/http"
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
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
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
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}
