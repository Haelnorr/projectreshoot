package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Wraps the http.ResponseWriter, adding a statusCode field
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

// Extends WriteHeader to the ResponseWriter to add the status code
func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Middleware to add logs to console with details of the request
func Logging(logger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		logger.Info().
			Int("status", wrapped.statusCode).
			Str("method", r.Method).
			Str("resource", r.URL.Path).
			Dur("time_elapsed", time.Since(start)).Msg("Served")
	})
}
