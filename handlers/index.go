package handlers

import (
	"net/http"
	"projectreshoot/view/page"
)

func HandleRoot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				page.Error("404", "Page not found", "The page or resource you have requested does not exist").Render(r.Context(), w)
				return
			}
			page.Index().Render(r.Context(), w)
		},
	)
}
