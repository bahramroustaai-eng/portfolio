package transport

import "net/http"

func NewRouter(h *UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/user", h.CreateUser)
	mux.HandleFunc("POST /api/v1/login", h.Login)
	mux.Handle("/", http.FileServer(http.Dir("web/static")))
	return Logging(mux)
}
