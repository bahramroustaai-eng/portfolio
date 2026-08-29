package transport

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Module interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(mods ...Module) http.Handler {
	mux := http.NewServeMux()

	for _, m := range mods {
		m.RegisterRoutes(mux)
	}
	mux.Handle("/", http.FileServer(http.Dir("web/static")))
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	return Logging(mux)
}
