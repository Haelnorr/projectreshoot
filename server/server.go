package server

import (
	"database/sql"
	"net/http"

	"projectreshoot/middleware"

	"github.com/rs/zerolog"
)

// Returns a new http.Handler with all the routes and middleware added
func NewServer(
	config *Config,
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
	handler = middleware.Logging(logger, handler)
	return handler
}
