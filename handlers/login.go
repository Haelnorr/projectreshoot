package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/view/component/form"
	"projectreshoot/view/page"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Validates the username matches a user in the database and the password
// is correct. Returns the corresponding user
func validateLogin(conn *sql.DB, r *http.Request) (db.User, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")
	user, err := db.GetUserFromUsername(conn, formUsername)
	if err != nil {
		return db.User{}, errors.Wrap(err, "db.GetUserFromUsername")
	}

	err = user.CheckPassword(formPassword)
	if err != nil {
		return db.User{}, errors.New("Username or password incorrect")
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
	logger *zerolog.Logger,
	conn *sql.DB,
	secretKey string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			user, err := validateLogin(conn, r)
			if err != nil {
				form.LoginForm(err.Error()).Render(r.Context(), w)
				if err.Error() != "Username or password incorrect" {
					logger.Warn().Caller().Err(err).Msg("Login request failed")
				}
				return
			}

			// TODO: login success, use the userID to set the session
			rememberMe := checkRememberMe(r)
			fmt.Printf(
				"Login success, user: %v - remember me?: %t\n",
				user.Username,
				rememberMe,
			)

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
