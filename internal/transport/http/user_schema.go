package transport

import "time"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
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

type GetUsersResponse struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}
