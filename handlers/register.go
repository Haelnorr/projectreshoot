package handlers

import (
	"database/sql"
	"net/http"

	"projectreshoot/config"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/view/component/form"
	"projectreshoot/view/page"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func validateRegistration(conn *sql.DB, r *http.Request) (*db.User, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")
	formConfirmPassword := r.FormValue("confirm-password")
	unique, err := db.CheckUsernameUnique(conn, formUsername)
	if err != nil {
		return nil, errors.Wrap(err, "db.CheckUsernameUnique")
	}
	if !unique {
		return nil, errors.New("Username is taken")
	}
	if formPassword != formConfirmPassword {
		return nil, errors.New("Passwords do not match")
	}
	user, err := db.CreateNewUser(conn, formUsername, formPassword)
	if err != nil {
		return nil, errors.Wrap(err, "db.CreateNewUser")
	}

	return user, nil
}

func HandleRegisterRequest(
	config *config.Config,
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			user, err := validateRegistration(conn, r)
			if err != nil {
				if err.Error() != "Username is taken" &&
					err.Error() != "Passwords do not match" {
					logger.Warn().Caller().Err(err).Msg("Registration request failed")
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					form.RegisterForm(err.Error()).Render(r.Context(), w)
				}
				return
			}

			rememberMe := checkRememberMe(r)
			err = cookies.SetTokenCookies(w, r, config, user, true, rememberMe)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Warn().Caller().Err(err).Msg("Failed to set token cookies")
			}

			pageFrom := cookies.CheckPageFrom(w, r)
			w.Header().Set("HX-Redirect", pageFrom)
		},
	)
}

// Handles a request to view the login page. Will attempt to set "pagefrom"
// cookie so a successful login can redirect the user to the page they came
func HandleRegisterPage(trustedHost string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookies.SetPageFrom(w, r, trustedHost)
			page.Register().Render(r.Context(), w)
		},
	)
}
