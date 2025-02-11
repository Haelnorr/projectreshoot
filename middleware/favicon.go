package middleware

import (
	"net/http"
)

func Favicon(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/favicon.ico" {
				http.ServeFile(w, r, "static/favicon.ico")
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}
