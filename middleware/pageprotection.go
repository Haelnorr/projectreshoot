package middleware

import (
	"net/http"
	"projectreshoot/contexts"
	"projectreshoot/view/page"
)

// Checks if the user is set in the context and shows 401 page if not logged in
func RequiresLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contexts.GetUser(r.Context())
		if user == nil {
			page.Error(
				"401",
				"Unauthorized",
				"Please login to view this page",
			).Render(r.Context(), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Checks if the user is set in the context and redirects them to profile if
// they are logged in
func RequiresLogout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contexts.GetUser(r.Context())
		if user != nil {
			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
