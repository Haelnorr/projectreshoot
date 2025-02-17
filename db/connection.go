package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

// Wraps the database handle, providing a mutex to safely manage transactions
type SafeConn struct {
	db  *sql.DB
	mux sync.RWMutex
}

// Extends sql.Tx for use with SafeConn
type SafeTX struct {
	tx *sql.Tx
	sc *SafeConn
}

// Starts a new transaction, waiting up to 10 seconds if the database is locked
func (conn *SafeConn) Begin(ctx context.Context) (*SafeTX, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	lockAcquired := make(chan struct{})
	go func() {
		conn.mux.RLock()
		close(lockAcquired)
	}()

	select {
	case <-lockAcquired:
		tx, err := conn.db.BeginTx(ctx, nil)
		if err != nil {
			conn.mux.RUnlock()
			return nil, err
		}
		return &SafeTX{tx: tx, sc: conn}, nil
	case <-ctx.Done():
		return nil, errors.New("Transaction time out due to database lock")
	}
}

// Query the database inside the transaction
func (stx *SafeTX) Query(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	if stx.tx == nil {
		return nil, errors.New("Cannot query without a transaction")
	}
	return stx.tx.QueryContext(ctx, query, args...)
}

// Exec a statement on the database inside the transaction
func (stx *SafeTX) Exec(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	if stx.tx == nil {
		return nil, errors.New("Cannot exec without a transaction")
	}
	return stx.tx.ExecContext(ctx, query, args...)
}

// Commit commits the transaction and releases the lock.
func (stx *SafeTX) Commit() error {
	if stx.tx == nil {
		return errors.New("Cannot commit without a transaction")
	}
	err := stx.tx.Commit()
	stx.tx = nil

	stx.releaseLock()
	return err
}

// Rollback aborts the transaction.
func (stx *SafeTX) Rollback() error {
	if stx.tx == nil {
		return errors.New("Cannot rollback without a transaction")
	}
	err := stx.tx.Rollback()
	stx.tx = nil
	stx.releaseLock()
	return err
}

// Release the read lock for the transaction
func (stx *SafeTX) releaseLock() {
	if stx.sc != nil {
		stx.sc.mux.RUnlock()
	}
}

// Pause blocks new transactions for a backup.
func (conn *SafeConn) Pause() {
	conn.mux.Lock() // Blocks all new transactions
}

// Resume allows transactions to proceed.
func (conn *SafeConn) Resume() {
	conn.mux.Unlock()
}

// Returns a database connection handle for the Turso DB
func ConnectToDatabase(dbName string) (*sql.DB, error) {
	file := fmt.Sprintf("file:%s.db", dbName)
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}
	return db, nil
}
