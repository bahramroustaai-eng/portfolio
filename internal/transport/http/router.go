package transport

import "net/http"

type Module interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(mods ...Module) http.Handler {
	mux := http.NewServeMux()

	for _, m := range mods {
		m.RegisterRoutes(mux)
	}
	mux.Handle("/", http.FileServer(http.Dir("web/static")))
	return Logging(mux)
}
