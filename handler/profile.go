package handler

import (
	"net/http"
	"projectreshoot/view/page"
)

func ProfilePage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			page.Profile().Render(r.Context(), w)
		},
	)
}
