package middleware

import (
	"net/http"
	"strings"
)

var excludedFiles = map[string]bool{
	"/static/css/output.css": true,
}

// Checks is path requested if for an excluded file and returns the file
// instead of passing the request onto the next middleware
func ExcludedFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if excludedFiles[r.URL.Path] {
				filePath := strings.TrimPrefix(r.URL.Path, "/")
				http.ServeFile(w, r, filePath)
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}
