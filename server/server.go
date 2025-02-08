package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"projectreshoot/middleware"

	"github.com/joho/godotenv"
)

type Config struct {
	TrustedHost string
	Host        string
	Port        string
}

func GetConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found.")
	}

	config := &Config{
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
		TrustedHost: os.Getenv("TRUSTED_HOST"),
	}
	if config.Host == "" {
		return nil, errors.New("Envar not set: HOST")
	}
	if config.Port == "" {
		return nil, errors.New("Envar not set: PORT")
	}
	if config.TrustedHost == "" {
		return nil, errors.New("Envar not set: TRUSTED_HOST")
	}

	return config, nil
}

func NewServer(config *Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
	)
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	return handler
}
