package handlers

import (
	"github.com/a-h/templ"
	"net/http"
)

func HandlePage(Page templ.Component) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Page.Render(r.Context(), w)
		},
	)
}
