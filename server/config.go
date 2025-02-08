package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host               string // Host to listen on
	Port               string // Port to listen on
	TrustedHost        string // Domain/Hostname to accept as trusted
	TursoDBName        string // DB Name for Turso DB/Branch
	TursoToken         string // Bearer token for Turso DB/Branch
	SecretKey          string // Secret key for signing tokens
	AccessTokenExpiry  int64  // Access token expiry in minutes
	RefreshTokenExpiry int64  // Refresh token expiry in minutes
	TokenFreshTime     int64  // Time for tokens to stay fresh in minutes
}

// Load the application configuration and get a pointer to the Config object
func GetConfig(args []string) (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found.")
	}
	var port string

	if args[0] != "" {
		port = args[0]
	} else {
		port = GetEnvDefault("PORT", "3333")
	}

	config := &Config{
		Host:               GetEnvDefault("HOST", "127.0.0.1"),
		Port:               port,
		TrustedHost:        os.Getenv("TRUSTED_HOST"),
		TursoDBName:        os.Getenv("TURSO_DB_NAME"),
		TursoToken:         os.Getenv("TURSO_AUTH_TOKEN"),
		SecretKey:          os.Getenv("SECRET_KEY"),
		AccessTokenExpiry:  GetEnvInt64("ACCESS_TOKEN_EXPIRY", 5),
		RefreshTokenExpiry: GetEnvInt64("REFRESH_TOKEN_EXPIRY", 1440), // defaults to 1 day
		TokenFreshTime:     GetEnvInt64("TOKEN_FRESH_TIME", 5),
	}

	if config.TrustedHost == "" {
		return nil, errors.New("Envar not set: TRUSTED_HOST")
	}
	if config.TursoDBName == "" {
		return nil, errors.New("Envar not set: TURSO_DB_NAME")
	}
	if config.TursoToken == "" {
		return nil, errors.New("Envar not set: TURSO_AUTH_TOKEN")
	}
	if config.SecretKey == "" {
		return nil, errors.New("Envar not set: SECRET_KEY")
	}

	return config, nil
}
