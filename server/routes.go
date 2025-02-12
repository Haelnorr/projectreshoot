package server

import (
	"database/sql"
	"net/http"

	"projectreshoot/config"
	"projectreshoot/handlers"
	"projectreshoot/view/page"

	"github.com/rs/zerolog"
)

// Add all the handled routes to the mux
func addRoutes(
	mux *http.ServeMux,
	logger *zerolog.Logger,
	config *config.Config,
	conn *sql.DB,
) {
	// Health check
	mux.HandleFunc("GET /healthz", func(http.ResponseWriter, *http.Request) {})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", handlers.HandleStatic()))

	// Index page and unhandled catchall (404)
	mux.Handle("GET /", handlers.HandleRoot())

	// Static content, unprotected pages
	mux.Handle("GET /about", handlers.HandlePage(page.About()))

	// Login page and handlers
	mux.Handle("GET /login", handlers.HandleLoginPage(config.TrustedHost))
	mux.Handle("POST /login", handlers.HandleLoginRequest(
		config,
		logger,
		conn,
		config.SecretKey,
	))

	// Profile page
}
