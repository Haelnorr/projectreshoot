package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"projectreshoot/config"
	"projectreshoot/db"
	"projectreshoot/logging"
	"projectreshoot/server"

	"github.com/pkg/errors"
)

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

	conn, err := db.ConnectToDatabase(&config.TursoDBName, &config.TursoToken)
	if err != nil {
		return errors.Wrap(err, "db.ConnectToDatabase")
	}
	defer conn.Close()

	srv := server.NewServer(config, logger, conn)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	// Runs function for testing in dev if --test flag true
	if args["test"] == "true" {
		test(config, logger, conn, httpServer)
		return nil
	}

	go func() {
		fmt.Fprintf(w, "Listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	fmt.Fprintln(w, "Shutting down")
	return nil
}

//go:embed static/*
var static embed.FS

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
