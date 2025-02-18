package handlers

import (
	"context"
	"net/http"
	"time"

	"projectreshoot/config"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/view/component/form"
	"projectreshoot/view/page"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Validates the username matches a user in the database and the password
// is correct. Returns the corresponding user
func validateLogin(
	ctx context.Context,
	tx *db.SafeTX,
	r *http.Request,
) (*db.User, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")
	user, err := db.GetUserFromUsername(ctx, tx, formUsername)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetUserFromUsername")
	}

	err = user.CheckPassword(formPassword)
	if err != nil {
		return nil, errors.New("Username or password incorrect")
	}
	return user, nil
}

// Returns result of the "Remember me?" checkbox as a boolean
func checkRememberMe(r *http.Request) bool {
	rememberMe := r.FormValue("remember-me")
	if rememberMe == "on" {
		return true
	} else {
		return false
	}
}

// Handles an attempted login request. On success will return a HTMX redirect
// and on fail will return the login form again, passing the error to the
// template for user feedback
func HandleLoginRequest(
	config *config.Config,
	logger *zerolog.Logger,
	conn *db.SafeConn,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()

			// Start the transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to set token cookies")
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			r.ParseForm()
			user, err := validateLogin(ctx, tx, r)
			if err != nil {
				tx.Rollback()
				if err.Error() != "Username or password incorrect" {
					logger.Warn().Caller().Err(err).Msg("Login request failed")
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					form.LoginForm(err.Error()).Render(r.Context(), w)
				}
				return
			}

			rememberMe := checkRememberMe(r)
			err = cookies.SetTokenCookies(w, r, config, user, true, rememberMe)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				logger.Warn().Caller().Err(err).Msg("Failed to set token cookies")
				return
			}

			tx.Commit()
			pageFrom := cookies.CheckPageFrom(w, r)
			w.Header().Set("HX-Redirect", pageFrom)
		},
	)
}

// Handles a request to view the login page. Will attempt to set "pagefrom"
// cookie so a successful login can redirect the user to the page they came
func HandleLoginPage(trustedHost string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookies.SetPageFrom(w, r, trustedHost)
			page.Login().Render(r.Context(), w)
		},
	)
}
