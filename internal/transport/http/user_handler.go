package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"portfolio/internal/user"
	"time"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (in CreateUserRequest) Validate() error {
	if len(in.Username) < 3 || len(in.Username) > 32 {
		return fmt.Errorf("%w: username must be between 3 and 32 characters", user.ErrInvalidInput)
	}
	if len(in.Password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", user.ErrInvalidInput)
	}
	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.svc.CreateUser(r.Context(), user.CreateInput{
		UserName: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateUserResponse{
		ID:        created.ID,
		Username:  created.UserName,
		CreatedAt: created.CreatedAt,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	writeJSON(w, statusCode, map[string]string{"error": msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, user.ErrUserNameConflict):
		writeError(w, http.StatusBadRequest, "username already exists")
	default:
		slog.Error("unhandled error", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
