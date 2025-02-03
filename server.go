package main

import (
	"net/http"
	"projectreshoot/middleware"
)

type Config struct {
	Host string
	Port string
}

func NewServer(
	config *Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		*config,
	)
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	return handler
}
