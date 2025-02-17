package handlers

import (
	"context"
	"net/http"
	"time"

	"projectreshoot/db"
	"projectreshoot/view/page"

	"github.com/rs/zerolog"
)

// A helper function to create a transaction with a cancellable context.
func WithTransaction(
	w http.ResponseWriter,
	r *http.Request,
	logger *zerolog.Logger,
	conn *db.SafeConn,
	handler func(
		ctx context.Context,
		tx *db.SafeTX,
		w http.ResponseWriter,
		r *http.Request,
	),
) {
	// Create a cancellable context from the request context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Start the transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		tx.Rollback()
		logger.Warn().Err(err).Msg("Request failed to start a transaction")
		w.WriteHeader(http.StatusServiceUnavailable)
		page.Error(
			"503",
			http.StatusText(503),
			"This service is currently unavailable. It could be down for maintenance").
			Render(r.Context(), w)
		return
	}

	// Pass the context and transaction to the handler
	handler(ctx, tx, w, r)
}
