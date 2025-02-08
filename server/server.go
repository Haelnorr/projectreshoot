package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"projectreshoot/middleware"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	TrustedHost string
	TursoURL    string
	TursoToken  string
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
		TursoURL:    os.Getenv("TURSO_DATABASE_URL"),
		TursoToken:  os.Getenv("TURSO_AUTH_TOKEN"),
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
	if config.TursoURL == "" {
		return nil, errors.New("Envar not set: TURSO_DATABASE_URL")
	}
	if config.TursoToken == "" {
		return nil, errors.New("Envar not set: TURSO_AUTH_TOKEN")
	}

	return config, nil
}

func NewServer(config *Config, conn *sql.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
		conn,
	)
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	return handler
}
