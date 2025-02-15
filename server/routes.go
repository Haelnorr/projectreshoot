package server

import (
	"database/sql"
	"net/http"

	"projectreshoot/config"
	"projectreshoot/handlers"
	"projectreshoot/middleware"
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
	mux.Handle("GET /login",
		middleware.RequiresLogout(
			handlers.HandleLoginPage(config.TrustedHost),
		))
	mux.Handle("POST /login",
		middleware.RequiresLogout(
			handlers.HandleLoginRequest(
				config,
				logger,
				conn,
			)))

	// Register page and handlers
	mux.Handle("GET /register",
		middleware.RequiresLogout(
			handlers.HandleRegisterPage(config.TrustedHost),
		))
	mux.Handle("POST /register",
		middleware.RequiresLogout(
			handlers.HandleRegisterRequest(
				config,
				logger,
				conn,
			)))

	// Logout
	mux.Handle("POST /logout", handlers.HandleLogout(config, logger, conn))

	// Profile page
	mux.Handle("GET /profile",
		middleware.RequiresLogin(
			handlers.HandleProfilePage(),
		))

	// Account page
	mux.Handle("GET /account",
		middleware.RequiresLogin(
			handlers.HandleAccountPage(),
		))
	mux.Handle("POST /account-select-page",
		middleware.RequiresLogin(
			handlers.HandleAccountSubpage(),
		))
	mux.Handle("POST /change-username",
		middleware.RequiresLogin(
			middleware.RequiresFresh(
				handlers.HandleChangeUsername(logger, conn),
			),
		))
	mux.Handle("POST /reauthenticate",
		middleware.RequiresLogin(
			handlers.HandleReauthenticate(logger, config, conn),
		))
}
