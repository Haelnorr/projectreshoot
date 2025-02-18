package handlers

import (
	"context"
	"net/http"
	"time"

	"projectreshoot/config"
	"projectreshoot/contexts"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/jwt"
	"projectreshoot/view/component/form"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Get the tokens from the request
func getTokens(
	config *config.Config,
	ctx context.Context,
	tx *db.SafeTX,
	r *http.Request,
) (*jwt.AccessToken, *jwt.RefreshToken, error) {
	// get the existing tokens from the cookies
	atStr, rtStr := cookies.GetTokenStrings(r)
	aT, err := jwt.ParseAccessToken(config, ctx, tx, atStr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "jwt.ParseAccessToken")
	}
	rT, err := jwt.ParseRefreshToken(config, ctx, tx, rtStr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "jwt.ParseRefreshToken")
	}
	return aT, rT, nil
}

// Revoke the given token pair
func revokeTokenPair(
	ctx context.Context,
	tx *db.SafeTX,
	aT *jwt.AccessToken,
	rT *jwt.RefreshToken,
) error {
	err := jwt.RevokeToken(ctx, tx, aT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	err = jwt.RevokeToken(ctx, tx, rT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	return nil
}

// Issue new tokens for the user, invalidating the old ones
func refreshTokens(
	config *config.Config,
	ctx context.Context,
	tx *db.SafeTX,
	w http.ResponseWriter,
	r *http.Request,
) error {
	aT, rT, err := getTokens(config, ctx, tx, r)
	if err != nil {
		return errors.Wrap(err, "getTokens")
	}
	rememberMe := map[string]bool{
		"session": false,
		"exp":     true,
	}[aT.TTL]
	// issue new tokens for the user
	user := contexts.GetUser(r.Context())
	err = cookies.SetTokenCookies(w, r, config, user.User, true, rememberMe)
	if err != nil {
		return errors.Wrap(err, "cookies.SetTokenCookies")
	}
	err = revokeTokenPair(ctx, tx, aT, rT)
	if err != nil {
		return errors.Wrap(err, "revokeTokenPair")
	}

	return nil
}

// Validate the provided password
func validatePassword(
	r *http.Request,
) error {
	r.ParseForm()
	password := r.FormValue("password")
	user := contexts.GetUser(r.Context())
	err := user.CheckPassword(password)
	if err != nil {
		return errors.Wrap(err, "user.CheckPassword")
	}
	return nil
}

// Handle request to reauthenticate (i.e. make token fresh again)
func HandleReauthenticate(
	logger *zerolog.Logger,
	config *config.Config,
	conn *db.SafeConn,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()

			// Start the transaction
			tx, err := conn.Begin(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to refresh user tokens")
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			err = validatePassword(r)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(445)
				form.ConfirmPassword("Incorrect password").Render(r.Context(), w)
				return
			}
			err = refreshTokens(config, ctx, tx, w, r)
			if err != nil {
				tx.Rollback()
				logger.Error().Err(err).Msg("Failed to refresh user tokens")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tx.Commit()
			w.WriteHeader(http.StatusOK)
		},
	)
}
