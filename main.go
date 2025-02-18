package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"projectreshoot/config"
	"projectreshoot/db"
	"projectreshoot/logging"
	"projectreshoot/server"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:embed static/*
var embeddedStatic embed.FS

// Gets the static files
func getStaticFiles(logger *zerolog.Logger) (http.FileSystem, error) {
	if _, err := os.Stat("static"); err == nil {
		// Use actual filesystem in development
		logger.Debug().Msg("Using filesystem for static files")
		return http.Dir("static"), nil
	} else {
		// Use embedded filesystem in production
		logger.Debug().Msg("Using embedded static files")
		subFS, err := fs.Sub(embeddedStatic, "static")
		if err != nil {
			return nil, errors.Wrap(err, "fs.Sub")
		}
		return http.FS(subFS), nil
	}
}

var maint uint32 // atomic: 1 if in maintenance mode

func handleMaintSignals(
	conn *db.SafeConn,
	srv *http.Server,
	logger *zerolog.Logger,
	config *config.Config,
) {
	ch := make(chan os.Signal, 1)
	srv.RegisterOnShutdown(func() {
		close(ch)
	})
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGUSR1:
				if atomic.LoadUint32(&maint) != 1 {
					atomic.StoreUint32(&maint, 1)
					log := logger.With().Logger().Output(os.Stdout)
					log.Info().Msg("Signal received: Starting maintenance")
					log.Info().Msg("Attempting to acquire database lock")
					conn.Pause(config.DBLockTimeout * time.Second)
				}
			case syscall.SIGUSR2:
				if atomic.LoadUint32(&maint) != 0 {
					log := logger.With().Logger().Output(os.Stdout)
					log.Info().Msg("Signal received: Maintenance over")
					log.Info().Msg("Releasing database lock")
					conn.Resume()
					atomic.StoreUint32(&maint, 0)
				}
			}
		}
	}()
	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
}

// Initializes and runs the server
func run(ctx context.Context, w io.Writer, args map[string]string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := config.GetConfig(args)
	if err != nil {
		return errors.Wrap(err, "server.GetConfig")
	}

	var logfile *os.File = nil
	if config.LogOutput == "both" || config.LogOutput == "file" {
		logfile, err = logging.GetLogFile(config.LogDir)
		if err != nil {
			return errors.Wrap(err, "logging.GetLogFile")
		}
		defer logfile.Close()
	}

	var consoleWriter io.Writer
	if config.LogOutput == "both" || config.LogOutput == "console" {
		consoleWriter = w
	}

	logger, err := logging.GetLogger(
		config.LogLevel,
		consoleWriter,
		logfile,
		config.LogDir,
	)
	if err != nil {
		return errors.Wrap(err, "logging.GetLogger")
	}

	conn, err := db.ConnectToDatabase(config.DBName, logger)
	if err != nil {
		return errors.Wrap(err, "db.ConnectToDatabase")
	}
	defer conn.Close()

	staticFS, err := getStaticFiles(logger)
	if err != nil {
		return errors.Wrap(err, "getStaticFiles")
	}

	srv := server.NewServer(config, logger, conn, &staticFS, &maint)
	httpServer := &http.Server{
		Addr:              net.JoinHostPort(config.Host, config.Port),
		Handler:           srv,
		ReadHeaderTimeout: config.ReadHeaderTimeout * time.Second,
		WriteTimeout:      config.WriteTimeout * time.Second,
		IdleTimeout:       config.IdleTimeout * time.Second,
	}

	// Runs function for testing in dev if --test flag true
	if args["test"] == "true" {
		test(config, logger, conn, httpServer)
		return nil
	}

	// Setups a channel to listen for os.Signal
	handleMaintSignals(conn, httpServer, logger, config)

	// Runs the http server
	go func() {
		logger.Info().Str("address", httpServer.Addr).Msg("Listening for requests")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("Error listening and serving")
		}
	}()

	// Handles graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("Error shutting down server")
		}
	}()
	wg.Wait()
	logger.Info().Msg("Shutting down")
	return nil
}

// Start of runtime. Parse commandline arguments & flags, Initializes context
// and starts the server
func main() {
	// Parse commandline args
	host := flag.String("host", "", "Override host to listen on")
	port := flag.String("port", "", "Override port to listen on")
	test := flag.Bool("test", false, "Run test function instead of main program")
	loglevel := flag.String("loglevel", "", "Set log level")
	logoutput := flag.String("logoutput", "", "Set log destination (file, console or both)")
	flag.Parse()

	// Map the args for easy access
	args := map[string]string{
		"host":      *host,
		"port":      *port,
		"test":      strconv.FormatBool(*test),
		"loglevel":  *loglevel,
		"logoutput": *logoutput,
	}

	// Start the server
	ctx := context.Background()
	if err := run(ctx, os.Stdout, args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
