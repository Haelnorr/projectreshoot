package server

import (
	"database/sql"
	"net/http"

	"projectreshoot/config"
	"projectreshoot/middleware"

	"github.com/rs/zerolog"
)

// Returns a new http.Handler with all the routes and middleware added
func NewServer(
	config *config.Config,
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		config,
		conn,
	)
	var handler http.Handler = mux
	// Add middleware here, must be added in reverse order of execution
	// i.e. First in list will get executed last during the request handling
	handler = middleware.Logging(logger, handler)
	handler = middleware.Authentication(logger, config, conn, handler)

	// Serve the favicon and exluded files before any middleware is added
	handler = middleware.ExcludedFiles(handler)
	handler = middleware.Favicon(handler)

	// Start the timer for the request chain so logger can have accurate info
	handler = middleware.StartTimer(handler)
	return handler
}
