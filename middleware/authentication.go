package middleware

import (
	"database/sql"
	"net/http"

	"projectreshoot/config"
	"projectreshoot/contexts"
	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/jwt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Attempt to use a valid refresh token to generate a new token pair
func refreshAuthTokens(
	config *config.Config,
	conn *sql.DB,
	w http.ResponseWriter,
	req *http.Request,
	ref *jwt.RefreshToken,
) (*db.User, error) {
	user, err := ref.GetUser(conn)
	if err != nil {
		return nil, errors.Wrap(err, "rT.GetUser")
	}

	rememberMe := map[string]bool{
		"session": false,
		"exp":     true,
	}[ref.TTL]

	// Set fresh to true because new tokens coming from refresh request
	err = cookies.SetTokenCookies(w, req, config, user, false, rememberMe)
	if err != nil {
		return nil, errors.Wrap(err, "cookies.SetTokenCookies")
	}
	// New tokens sent, revoke the used refresh token
	err = jwt.RevokeToken(conn, ref)
	if err != nil {
		return nil, errors.Wrap(err, "jwt.RevokeToken")
	}
	// Return the authorized user
	return user, nil
}

// Check the cookies for token strings and attempt to authenticate them
func getAuthenticatedUser(
	config *config.Config,
	conn *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
) (*db.User, error) {
	// Get token strings from cookies
	atStr, rtStr := cookies.GetTokenStrings(r)
	// Attempt to parse the access token
	aT, err := jwt.ParseAccessToken(config, conn, atStr)
	if err != nil {
		// Access token invalid, attempt to parse refresh token
		rT, err := jwt.ParseRefreshToken(config, conn, rtStr)
		if err != nil {
			return nil, errors.Wrap(err, "jwt.ParseRefreshToken")
		}
		// Refresh token valid, attempt to get a new token pair
		user, err := refreshAuthTokens(config, conn, w, r, rT)
		if err != nil {
			return nil, errors.Wrap(err, "refreshAuthTokens")
		}
		// New token pair sent, return the authorized user
		return user, nil
	}
	// Access token valid
	user, err := aT.GetUser(conn)
	if err != nil {
		return nil, errors.Wrap(err, "rT.GetUser")
	}
	return user, nil
}

// Attempt to authenticate the user and add their account details
// to the request context
func Authentication(
	logger *zerolog.Logger,
	config *config.Config,
	conn *sql.DB,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getAuthenticatedUser(config, conn, w, r)
		if err != nil {
			// User auth failed, delete the cookies to avoid repeat requests
			cookies.DeleteCookie(w, "access", "/")
			cookies.DeleteCookie(w, "refresh", "/")
			logger.Debug().
				Str("remote_addr", r.RemoteAddr).
				Err(err).
				Msg("Failed to authenticate user")
		}
		ctx := contexts.SetUser(r.Context(), user)
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)
	})
}
