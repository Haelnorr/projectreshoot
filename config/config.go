package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"projectreshoot/logging"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Host               string        // Host to listen on
	Port               string        // Port to listen on
	TrustedHost        string        // Domain/Hostname to accept as trusted
	SSL                bool          // Flag for SSL Mode
	ReadHeaderTimeout  time.Duration // Timeout for reading request headers in seconds
	WriteTimeout       time.Duration // Timeout for writing requests in seconds
	IdleTimeout        time.Duration // Timeout for idle connections in seconds
	DBName             string        // Filename of the db (doesnt include file extension)
	SecretKey          string        // Secret key for signing tokens
	AccessTokenExpiry  int64         // Access token expiry in minutes
	RefreshTokenExpiry int64         // Refresh token expiry in minutes
	TokenFreshTime     int64         // Time for tokens to stay fresh in minutes
	LogLevel           zerolog.Level // Log level for global logging. Defaults to info
	LogOutput          string        // "file", "console", or "both". Defaults to console
	LogDir             string        // Path to create log files
}

// Load the application configuration and get a pointer to the Config object
func GetConfig(args map[string]string) (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	var (
		host      string
		port      string
		logLevel  zerolog.Level
		logOutput string
		valid     bool
	)

	if args["host"] != "" {
		host = args["host"]
	} else {
		host = GetEnvDefault("HOST", "127.0.0.1")
	}
	if args["port"] != "" {
		port = args["port"]
	} else {
		port = GetEnvDefault("PORT", "3333")
	}
	if args["loglevel"] != "" {
		logLevel = logging.GetLogLevel(args["loglevel"])
	} else {
		logLevel = logging.GetLogLevel(GetEnvDefault("LOG_LEVEL", "info"))
	}
	if args["logoutput"] != "" {
		opts := map[string]string{
			"both":    "both",
			"file":    "file",
			"console": "console",
		}
		logOutput, valid = opts[args["logoutput"]]
		if !valid {
			logOutput = "console"
			fmt.Println(
				"Log output type was not parsed correctly. Defaulting to console only",
			)
		}
	} else {
		logOutput = GetEnvDefault("LOG_OUTPUT", "console")
	}
	if logOutput != "both" && logOutput != "console" && logOutput != "file" {
		logOutput = "console"
	}

	config := &Config{
		Host:               host,
		Port:               port,
		TrustedHost:        GetEnvDefault("TRUSTED_HOST", "127.0.0.1"),
		SSL:                GetEnvBool("SSL_MODE", false),
		ReadHeaderTimeout:  GetEnvDur("READ_HEADER_TIMEOUT", 2),
		WriteTimeout:       GetEnvDur("WRITE_TIMEOUT", 10),
		IdleTimeout:        GetEnvDur("IDLE_TIMEOUT", 120),
		DBName:             GetEnvDefault("DB_NAME", "projectreshoot"),
		SecretKey:          os.Getenv("SECRET_KEY"),
		AccessTokenExpiry:  GetEnvInt64("ACCESS_TOKEN_EXPIRY", 5),
		RefreshTokenExpiry: GetEnvInt64("REFRESH_TOKEN_EXPIRY", 1440), // defaults to 1 day
		TokenFreshTime:     GetEnvInt64("TOKEN_FRESH_TIME", 5),
		LogLevel:           logLevel,
		LogOutput:          logOutput,
		LogDir:             GetEnvDefault("LOG_DIR", ""),
	}

	if config.SecretKey == "" {
		return nil, errors.New("Envar not set: SECRET_KEY")
	}

	return config, nil
}
