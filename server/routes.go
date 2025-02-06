package server

import (
	"net/http"
	"projectreshoot/handlers"
	"projectreshoot/view/page"
)

func addRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", handlers.HandleStatic()))
	mux.Handle("GET /", handlers.HandleRoot())
	mux.Handle("GET /about", handlers.HandlePage(page.About()))
	mux.Handle("GET /login", handlers.HandlePage(page.Login()))
}
