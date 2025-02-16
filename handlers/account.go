package handlers

import (
	"database/sql"
	"net/http"

	"projectreshoot/contexts"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/view/component/account"
	"projectreshoot/view/page"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Renders the account page on the 'General' subpage
func HandleAccountPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("subpage")
			subpage := cookie.Value
			if err != nil {
				subpage = "General"
			}
			page.Account(subpage).Render(r.Context(), w)
		},
	)
}

// Handles a request to change the subpage for the Accou/accountnt page
func HandleAccountSubpage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			subpage := r.FormValue("subpage")
			cookies.SetCookie(w, "subpage", "/account", subpage, 300)
			account.AccountContainer(subpage).Render(r.Context(), w)
		},
	)
}

// Handles a request to change the users username
func HandleChangeUsername(
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			newUsername := r.FormValue("username")

			unique, err := db.CheckUsernameUnique(conn, newUsername)
			if err != nil {
				logger.Error().Err(err).Msg("Error updating username")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !unique {
				account.ChangeUsername("Username is taken", newUsername).
					Render(r.Context(), w)
				return
			}
			user := contexts.GetUser(r.Context())
			err = user.ChangeUsername(conn, newUsername)
			if err != nil {
				logger.Error().Err(err).Msg("Error updating username")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Refresh", "true")
		},
	)
}

// Handles a request to change the users bio
func HandleChangeBio(
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			newBio := r.FormValue("bio")
			leng := len([]rune(newBio))
			if leng > 128 {
				account.ChangeBio("Bio limited to 128 characters", newBio).
					Render(r.Context(), w)
				return
			}
			user := contexts.GetUser(r.Context())
			err := user.ChangeBio(conn, newBio)
			if err != nil {
				logger.Error().Err(err).Msg("Error updating bio")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Refresh", "true")
		},
	)
}
func validateChangePassword(conn *sql.DB, r *http.Request) (string, error) {
	r.ParseForm()
	formPassword := r.FormValue("password")
	formConfirmPassword := r.FormValue("confirm-password")
	if formPassword != formConfirmPassword {
		return "", errors.New("Passwords do not match")
	}
	if len(formPassword) > 72 {
		return "", errors.New("Password exceeds maximum length of 72 bytes")
	}
	return formPassword, nil
}

// Handles a request to change the users password
func HandleChangePassword(
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			newPass, err := validateChangePassword(conn, r)
			if err != nil {
				account.ChangePassword(err.Error()).Render(r.Context(), w)
				return
			}
			user := contexts.GetUser(r.Context())
			err = user.SetPassword(conn, newPass)
			if err != nil {
				logger.Error().Err(err).Msg("Error updating password")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Refresh", "true")
		},
	)
}
