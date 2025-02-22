package handler

import (
	"net/http"

	"projectreshoot/view/page"
)

// Handles responses to the / path. Also serves a 404 Page for paths that
// don't have explicit handlers
func Root() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				ErrorPage(http.StatusNotFound, w, r)
				return
			}
			page.Index().Render(r.Context(), w)
		},
	)
}
