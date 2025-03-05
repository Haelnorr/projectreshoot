package main

import (
	"net/http"
	"os"
	"os/signal"
	"projectreshoot/pkg/config"
	"projectreshoot/pkg/db"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

// Handle SIGUSR1 and SIGUSR2 syscalls to toggle maintenance mode
func handleMaintSignals(
	conn *db.SafeConn,
	config *config.Config,
	srv *http.Server,
	logger *zerolog.Logger,
) {
	logger.Debug().Msg("Starting signal listener")
	ch := make(chan os.Signal, 1)
	srv.RegisterOnShutdown(func() {
		logger.Debug().Msg("Shutting down signal listener")
		close(ch)
	})
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGUSR1:
				if atomic.LoadUint32(&maint) != 1 {
					atomic.StoreUint32(&maint, 1)
					logger.Info().Msg("Signal received: Starting maintenance")
					logger.Info().Msg("Attempting to acquire database lock")
					conn.Pause(config.DBLockTimeout * time.Second)
				}
			case syscall.SIGUSR2:
				if atomic.LoadUint32(&maint) != 0 {
					logger.Info().Msg("Signal received: Maintenance over")
					logger.Info().Msg("Releasing database lock")
					conn.Resume()
					atomic.StoreUint32(&maint, 0)
				}
			}
		}
	}()
	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
}
