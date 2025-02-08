package server

import (
	"net/http"
	"projectreshoot/handlers"
	"projectreshoot/view/page"
)

func addRoutes(
	mux *http.ServeMux,
) {
	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", handlers.HandleStatic()))

	// Index page
	mux.Handle("GET /", handlers.HandleRoot())

	// Static pages
	mux.Handle("GET /about", handlers.HandlePage(page.About()))

	// Login page and handlers
	mux.Handle("GET /login", handlers.HandleLoginPage())
	mux.Handle("POST /login", handlers.HandleLoginRequest())
}
