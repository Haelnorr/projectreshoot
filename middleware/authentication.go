package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// Take current request
// Get cookies from browser
// Parse the tokens
// Check if tokens blacklisted
// Trigger refresh if required
// Create context with state of user authorization
// Pass request on with context

func Authentication(logger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
