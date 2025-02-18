package handlers

import (
	"context"
	"net/http"
	"time"

	"projectreshoot/db"

	"github.com/rs/zerolog"
)

func removeme(
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
	onfail func(err error),
) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	// Start the transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		onfail(err)
		return
	}

	handler(ctx, tx, w, r)
}
