package main

import (
	"net/http"
	"projectreshoot/handlers"
)

func addRoutes(
	mux *http.ServeMux,
	config Config,
) {
	mux.Handle("GET /", handlers.HandleRoot())
}
