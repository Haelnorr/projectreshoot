package handlers

import (
	"html/template"
	"net/http"
)

func HandleRoot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			templ := template.Must(template.ParseFiles("templates/index.html"))
			templ.Execute(w, nil)
		},
	)
}
