package transport

import "net/http"

func NewRouter(h *UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/user", h.CreateUser)
	return Logging(mux)
}
