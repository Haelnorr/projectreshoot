package middleware

import (
	"net/http"
	"projectreshoot/contexts"
	"time"
)

func StartTimer(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := contexts.SetStart(r.Context(), start)
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		},
	)
}
