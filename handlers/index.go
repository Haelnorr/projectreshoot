package handlers

import (
	"net/http"

	"projectreshoot/view/page"
)

// Handles responses to the / path. Also serves a 404 Page for paths that
// don't have explicit handlers
func HandleRoot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				page.Error(
					"404",
					"Page not found",
					"The page or resource you have requested does not exist",
				).Render(r.Context(), w)
				return
			}
			page.Index().Render(r.Context(), w)
		},
	)
}
