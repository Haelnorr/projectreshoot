package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"projectreshoot/db"
	"projectreshoot/server"

	"github.com/pkg/errors"
)

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := server.GetConfig()
	if err != nil {
		return errors.Wrap(err, "server.GetConfig")
	}

	conn, err := db.ConnectToDatabase(&config.TursoURL, &config.TursoToken)
	if err != nil {
		return errors.Wrap(err, "db.ConnectToDatabase")
	}

	srv := server.NewServer(config, conn)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
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

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
