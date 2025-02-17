package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	_ "modernc.org/sqlite"
)

// Wraps the database handle, providing a mutex to safely manage transactions
type SafeConn struct {
	db  *sql.DB
	mux sync.RWMutex
}

func MakeSafe(db *sql.DB) *SafeConn {
	return &SafeConn{db: db}
}

// Extends sql.Tx for use with SafeConn
type SafeTX struct {
	tx *sql.Tx
	sc *SafeConn
}

// Starts a new transaction based on the current context. Will cancel if
// the context is closed/cancelled/done
func (conn *SafeConn) Begin(ctx context.Context) (*SafeTX, error) {
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

// Close the database connection
func (conn *SafeConn) Close() error {
	conn.mux.Lock()
	defer conn.mux.Unlock()
	return conn.db.Close()
}

// Returns a database connection handle for the DB
func OldConnectToDatabase(dbName string) (*sql.DB, error) {
	file := fmt.Sprintf("file:%s.db", dbName)
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	return db, nil
}

// Returns a database connection handle for the DB
func ConnectToDatabase(dbName string) (*SafeConn, error) {
	file := fmt.Sprintf("file:%s.db", dbName)
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	conn := MakeSafe(db)

	return conn, nil
}
