package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"projectreshoot/config"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/jwt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func revokeAccess(
	config *config.Config,
	ctx context.Context,
	tx *db.SafeTX,
	atStr string,
) error {
	aT, err := jwt.ParseAccessToken(config, ctx, tx, atStr)
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") ||
			strings.Contains(err.Error(), "Token has been revoked") {
			return nil // Token is expired, dont need to revoke it
		}
		return errors.Wrap(err, "jwt.ParseAccessToken")
	}
	err = jwt.RevokeToken(ctx, tx, aT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	return nil
}

func revokeRefresh(
	config *config.Config,
	ctx context.Context,
	tx *db.SafeTX,
	rtStr string,
) error {
	rT, err := jwt.ParseRefreshToken(config, ctx, tx, rtStr)
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") ||
			strings.Contains(err.Error(), "Token has been revoked") {
			return nil // Token is expired, dont need to revoke it
		}
		return errors.Wrap(err, "jwt.ParseRefreshToken")
	}
	err = jwt.RevokeToken(ctx, tx, rT)
	if err != nil {
		return errors.Wrap(err, "jwt.RevokeToken")
	}
	return nil
}

// Retrieve and revoke the user's tokens
func revokeTokens(
	config *config.Config,
	ctx context.Context,
	tx *db.SafeTX,
	r *http.Request,
) error {
	// get the tokens from the cookies
	atStr, rtStr := cookies.GetTokenStrings(r)
	// revoke the refresh token first as the access token expires quicker
	// only matters if there is an error revoking the tokens
	err := revokeRefresh(config, ctx, tx, rtStr)
	if err != nil {
		return errors.Wrap(err, "revokeRefresh")
	}
	err = revokeAccess(config, ctx, tx, atStr)
	if err != nil {
		return errors.Wrap(err, "revokeAccess")
	}
	return nil
}

// Handle a logout request
func Logout(
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
				logger.Warn().Err(err).Msg("Error occured on user logout")
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			err = revokeTokens(config, ctx, tx, r)
			if err != nil {
				tx.Rollback()
				logger.Error().Err(err).Msg("Error occured on user logout")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tx.Commit()
			cookies.DeleteCookie(w, "access", "/")
			cookies.DeleteCookie(w, "refresh", "/")
			w.Header().Set("HX-Redirect", "/login")
		},
	)
}
