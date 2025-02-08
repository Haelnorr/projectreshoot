package server

import (
	"database/sql"
	"net/http"

	"projectreshoot/middleware"
)

// Returns a new http.Handler with all the routes and middleware added
func NewServer(config *Config, conn *sql.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
		conn,
	)
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	return handler
}
