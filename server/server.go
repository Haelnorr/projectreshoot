package server

import (
	"net/http"
	"projectreshoot/middleware"
)

type Config struct {
	Host string
	Port string
}

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
	)
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	return handler
}
