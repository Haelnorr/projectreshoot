package handlers

import (
	"net/http"
	"projectreshoot/view/component/account"
	"projectreshoot/view/page"
)

func HandleAccountPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			page.Account("General").Render(r.Context(), w)
		},
	)
}

func HandleAccountSubpage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			subpage := r.FormValue("subpage")
			account.AccountContent(subpage).Render(r.Context(), w)
		},
	)
}
