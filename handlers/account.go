package handlers

import (
	"context"
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
			subpage := "General"
			if err == nil {
				subpage = cookie.Value
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
	conn *db.SafeConn,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WithTransaction(w, r, logger, conn,
				func(ctx context.Context, tx *db.SafeTX, w http.ResponseWriter, r *http.Request) {
					r.ParseForm()
					newUsername := r.FormValue("username")
					unique, err := db.CheckUsernameUnique(ctx, tx, newUsername)
					if err != nil {
						tx.Rollback()
						logger.Error().Err(err).Msg("Error updating username")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if !unique {
						tx.Rollback()
						account.ChangeUsername("Username is taken", newUsername).
							Render(r.Context(), w)
						return
					}
					user := contexts.GetUser(r.Context())
					err = user.ChangeUsername(ctx, tx, newUsername)
					if err != nil {
						tx.Rollback()
						logger.Error().Err(err).Msg("Error updating username")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					tx.Commit()
					w.Header().Set("HX-Refresh", "true")
				},
			)
		},
	)
}

// Handles a request to change the users bio
func HandleChangeBio(
	logger *zerolog.Logger,
	conn *db.SafeConn,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WithTransaction(w, r, logger, conn,
				func(ctx context.Context, tx *db.SafeTX, w http.ResponseWriter, r *http.Request) {
					r.ParseForm()
					newBio := r.FormValue("bio")
					leng := len([]rune(newBio))
					if leng > 128 {
						tx.Rollback()
						account.ChangeBio("Bio limited to 128 characters", newBio).
							Render(r.Context(), w)
						return
					}
					user := contexts.GetUser(r.Context())
					err := user.ChangeBio(ctx, tx, newBio)
					if err != nil {
						tx.Rollback()
						logger.Error().Err(err).Msg("Error updating bio")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					tx.Commit()
					w.Header().Set("HX-Refresh", "true")
				},
			)
		},
	)
}
func validateChangePassword(
	ctx context.Context,
	tx *db.SafeTX,
	r *http.Request,
) (string, error) {
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
	conn *db.SafeConn,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WithTransaction(w, r, logger, conn,
				func(ctx context.Context, tx *db.SafeTX, w http.ResponseWriter, r *http.Request) {
					newPass, err := validateChangePassword(ctx, tx, r)
					if err != nil {
						tx.Rollback()
						account.ChangePassword(err.Error()).Render(r.Context(), w)
						return
					}
					user := contexts.GetUser(r.Context())
					err = user.SetPassword(ctx, tx, newPass)
					if err != nil {
						tx.Rollback()
						logger.Error().Err(err).Msg("Error updating password")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					tx.Commit()
					w.Header().Set("HX-Refresh", "true")
				},
			)
		},
	)
}
