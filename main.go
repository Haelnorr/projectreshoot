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

	"projectreshoot/db"
	"projectreshoot/server"

	"github.com/pkg/errors"
)

// Initializes and runs the server
func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := server.GetConfig(args)
	if err != nil {
		return errors.Wrap(err, "server.GetConfig")
	}

	conn, err := db.ConnectToDatabase(&config.TursoDBName, &config.TursoToken)
	if err != nil {
		return errors.Wrap(err, "db.ConnectToDatabase")
	}

	srv := server.NewServer(config, conn)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	// TEST: runs function for testing in dev if --test flag true
	if args[1] == "true" {
		test(config, conn, httpServer)
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
	port := flag.String("port", "", "Override port")
	test := flag.Bool("test", false, "Run test function")
	flag.Parse()
	args := []string{*port, strconv.FormatBool(*test)}
	ctx := context.Background()
	if err := run(ctx, os.Stdout, args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
