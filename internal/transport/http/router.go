package transport

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const apiPrefix = "/api/v1"

type Module interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(mods ...Module) http.Handler {
	apiMux := http.NewServeMux()

	for _, m := range mods {
		m.RegisterRoutes(apiMux)
	}

	mux := http.NewServeMux()
	mux.Handle(
		apiPrefix+"/", http.StripPrefix(apiPrefix, apiMux),
	)
	mux.Handle("/", http.FileServer(http.Dir("web/static/")))
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	return Logging(mux)
}
