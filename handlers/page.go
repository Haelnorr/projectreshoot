package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

// Handler for static pages. Will render the given templ.Component to the
// http.ResponseWriter
func HandlePage(Page templ.Component) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Page.Render(r.Context(), w)
		},
	)
}
