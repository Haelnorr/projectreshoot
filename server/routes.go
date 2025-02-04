package server

import (
	"net/http"
	"projectreshoot/handlers"
)

func addRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", handlers.HandleStatic()))
	mux.Handle("GET /", handlers.HandleRoot())
}
