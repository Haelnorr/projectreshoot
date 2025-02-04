package handlers

import (
	"net/http"
	"projectreshoot/view/page"
)

func HandleRoot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			page.Index().Render(r.Context(), w)
		},
	)
}
