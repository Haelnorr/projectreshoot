package handlers

import (
	"database/sql"
	"net/http"

	"projectreshoot/contexts"
	"projectreshoot/db"
	"projectreshoot/view/component/account"
	"projectreshoot/view/page"

	"github.com/rs/zerolog"
)

// Renders the account page on the 'General' subpage
func HandleAccountPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			page.Account("General").Render(r.Context(), w)
		},
	)
}

// Handles a request to change the subpage for the Account page
func HandleAccountSubpage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			subpage := r.FormValue("subpage")
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
