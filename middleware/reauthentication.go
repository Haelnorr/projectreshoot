package middleware

import (
	"net/http"
	"projectreshoot/contexts"
	"time"
)

func RequiresFresh(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contexts.GetUser(r.Context())
		isFresh := time.Now().Before(time.Unix(user.Fresh, 0))
		if !isFresh {
			w.WriteHeader(444)
			return
		}
		next.ServeHTTP(w, r)
	})
}
