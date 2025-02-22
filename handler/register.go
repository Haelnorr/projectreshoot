package handler

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

func validateRegistration(
	ctx context.Context,
	tx *db.SafeTX,
	r *http.Request,
) (*db.User, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")
	formConfirmPassword := r.FormValue("confirm-password")
	unique, err := db.CheckUsernameUnique(ctx, tx, formUsername)
	if err != nil {
		return nil, errors.Wrap(err, "db.CheckUsernameUnique")
	}
	if !unique {
		return nil, errors.New("Username is taken")
	}
	if formPassword != formConfirmPassword {
		return nil, errors.New("Passwords do not match")
	}
	if len(formPassword) > 72 {
		return nil, errors.New("Password exceeds maximum length of 72 bytes")
	}
	user, err := db.CreateNewUser(ctx, tx, formUsername, formPassword)
	if err != nil {
		return nil, errors.Wrap(err, "db.CreateNewUser")
	}

	return user, nil
}

func RegisterRequest(
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
			user, err := validateRegistration(ctx, tx, r)
			if err != nil {
				tx.Rollback()
				if err.Error() != "Username is taken" &&
					err.Error() != "Passwords do not match" &&
					err.Error() != "Password exceeds maximum length of 72 bytes" {
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
func RegisterPage(trustedHost string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookies.SetPageFrom(w, r, trustedHost)
			page.Register().Render(r.Context(), w)
		},
	)
}
