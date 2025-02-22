package server

import (
	"net/http"

	"projectreshoot/config"
	"projectreshoot/db"
	"projectreshoot/middleware"

	"github.com/rs/zerolog"
)

// Returns a new http.Handler with all the routes and middleware added
func NewServer(
	config *config.Config,
	logger *zerolog.Logger,
	conn *db.SafeConn,
	staticFS *http.FileSystem,
	maint *uint32,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		config,
		conn,
		staticFS,
	)
	var handler http.Handler = mux
	// Add middleware here, must be added in reverse order of execution
	// i.e. First in list will get executed last during the request handling
	handler = middleware.Logging(logger, handler)
	handler = middleware.Authentication(logger, config, conn, handler, maint)

	// Gzip
	handler = middleware.Gzip(handler, config.GZIP)

	// Start the timer for the request chain so logger can have accurate info
	handler = middleware.StartTimer(handler)
	return handler
}
