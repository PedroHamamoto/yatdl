package user

import (
	"encoding/json"
	"log"
	"net/http"
	"yatdl/internal/http/httperr"
)

type Handler struct {
	service *Service
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httperr.Write(w, &httperr.BadRequest)
		log.Printf("failed to decode request: %v", err)
		return
	}

	err := h.service.Create(r.Context(), request.Email, request.Password)
	if err != nil {
		// TODO Handle errors properly
		log.Printf("failed to create user: %v", err)
		httperr.Write(w, &httperr.Unknown)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
