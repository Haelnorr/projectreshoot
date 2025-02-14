package handlers

import (
	"database/sql"
	"net/http"
	"projectreshoot/config"
	"projectreshoot/cookies"
	"projectreshoot/jwt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Retrieve and revoke the user's tokens
func revokeTokens(
	config *config.Config,
	conn *sql.DB,
	r *http.Request,
) error {
	// get the tokens from the cookies
	atStr, rtStr := cookies.GetTokenStrings(r)
	aT, err := jwt.ParseAccessToken(config, conn, atStr)
	if err != nil {
		return errors.Wrap(err, "jwt.ParseAccessToken")
	}
	rT, err := jwt.ParseRefreshToken(config, conn, rtStr)
	if err != nil {
		return errors.Wrap(err, "jwt.ParseRefreshToken")
	}
	// revoke the refresh token first as the access token expires quicker
	// only matters if there is an error revoking the tokens
	err = jwt.RevokeToken(conn, rT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	err = jwt.RevokeToken(conn, aT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	return nil
}

// Handle a logout request
func HandleLogout(
	config *config.Config,
	logger *zerolog.Logger,
	conn *sql.DB,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := revokeTokens(config, conn, r)
			if err != nil {
				logger.Error().Err(err).Msg("Error occured on user logout")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cookies.DeleteCookie(w, "access", "/")
			cookies.DeleteCookie(w, "refresh", "/")
			w.Header().Set("HX-Redirect", "/login")
		},
	)
}
