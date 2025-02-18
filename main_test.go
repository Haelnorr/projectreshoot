package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	args := map[string]string{}
	var stdout bytes.Buffer
	go run(ctx, &stdout, args)

	waitForReady(ctx, 10*time.Second, "http://localhost:3333/healthz")

	t.Run("SIGUSR1 puts database into global lock", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			expected := "Global database lock acquired"
			for {
				if strings.Contains(stdout.String(), expected) {
					done <- true
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		proc, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)
		proc.Signal(syscall.SIGUSR1)

		select {
		case <-done:
			t.Log("found")
		case <-time.After(250 * time.Millisecond):
			t.Errorf("Not found")
		}
	})

	t.Run("SIGUSR2 releases database global lock", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			expected := "Global database lock released"
			for {
				if strings.Contains(stdout.String(), expected) {
					done <- true
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		proc, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)
		proc.Signal(syscall.SIGUSR2)

		select {
		case <-done:
			t.Log("found")
		case <-time.After(250 * time.Millisecond):
			t.Errorf("Not found")
		}
	})
}

func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}
