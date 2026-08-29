package transport

import (
	"encoding/json"
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

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/user", h.CreateUser)
	mux.HandleFunc("POST /api/v1/login", h.Login)
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

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	Registers a new user with a username and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest	true	"user payload"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	created, err := h.svc.CreateUser(r.Context(), user.CreateInput{
		UserName: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	token, err := user.CreateToken(created.UserName, created.ID)
	if err != nil {
		slog.Error("create token", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		ID:          created.ID,
		Username:    created.UserName,
		AccessToken: token,
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          int32  `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

// Login godoc
//
//	@Summary		Login a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"user payload"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	loggedIn, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	token, err := user.CreateToken(loggedIn.UserName, loggedIn.ID)
	if err != nil {
		slog.Error("create token", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		ID:          loggedIn.ID,
		Username:    loggedIn.UserName,
		AccessToken: token,
	})
}
