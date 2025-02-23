package server

import (
	"net/http"

	"projectreshoot/config"
	"projectreshoot/db"
	"projectreshoot/handler"
	"projectreshoot/middleware"
	"projectreshoot/view/page"

	"github.com/rs/zerolog"
)

// Add all the handled routes to the mux
func addRoutes(
	mux *http.ServeMux,
	logger *zerolog.Logger,
	config *config.Config,
	conn *db.SafeConn,
	staticFS *http.FileSystem,
) {
	route := mux.Handle
	loggedIn := middleware.LoginReq
	loggedOut := middleware.LogoutReq
	fresh := middleware.FreshReq

	// Health check
	mux.HandleFunc("GET /healthz", func(http.ResponseWriter, *http.Request) {})

	// Static files
	route("GET /static/", http.StripPrefix("/static/", handler.StaticFS(staticFS)))

	// Index page and unhandled catchall (404)
	route("GET /", handler.Root())

	// Static content, unprotected pages
	route("GET /about", handler.HandlePage(page.About()))

	// Login page and handlers
	route("GET /login", loggedOut(handler.LoginPage(config.TrustedHost)))
	route("POST /login", loggedOut(handler.LoginRequest(config, logger, conn)))

	// Register page and handlers
	route("GET /register", loggedOut(handler.RegisterPage(config.TrustedHost)))
	route("POST /register", loggedOut(handler.RegisterRequest(config, logger, conn)))

	// Logout
	route("POST /logout", handler.Logout(config, logger, conn))

	// Reauthentication request
	route("POST /reauthenticate", loggedIn(handler.Reauthenticate(logger, config, conn)))

	// Profile page
	route("GET /profile", loggedIn(handler.ProfilePage()))

	// Account page
	route("GET /account", loggedIn(handler.AccountPage()))
	route("POST /account-select-page", loggedIn(handler.AccountSubpage()))
	route("POST /change-username", loggedIn(fresh(handler.ChangeUsername(logger, conn))))
	route("POST /change-bio", loggedIn(handler.ChangeBio(logger, conn)))
	route("POST /change-password", loggedIn(fresh(handler.ChangePassword(logger, conn))))

	// Movie page
	route("GET /movie/{movie_id}", handler.Movie(logger, config))
}
