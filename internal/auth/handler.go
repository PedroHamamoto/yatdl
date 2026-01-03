package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"yatdl/internal/http/httperr"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("failed to decode login request: %v", err)
		httperr.Write(w, &httperr.BadRequest)
		return
	}

	// TODO validate input

	response, err := h.service.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httperr.Write(w, &httperr.Unauthorized)
			return
		}

		log.Printf("Unknow Error happened in Login: %v", err)
		httperr.Write(w, &httperr.Unknown)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		httperr.Write(w, &httperr.Unknown)
		return
	}
}
