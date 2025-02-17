package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"

	_ "modernc.org/sqlite"
)

type SafeConn struct {
	db               *sql.DB
	readLockCount    int32
	globalLockStatus int32
}

func MakeSafe(db *sql.DB) *SafeConn {
	return &SafeConn{db: db}
}

// Extends sql.Tx for use with SafeConn
type SafeTX struct {
	tx *sql.Tx
	sc *SafeConn
}

func (conn *SafeConn) acquireGlobalLock() bool {
	if atomic.LoadInt32(&conn.readLockCount) > 0 || atomic.LoadInt32(&conn.globalLockStatus) == 1 {
		return false
	}
	atomic.StoreInt32(&conn.globalLockStatus, 1)
	return true
}

func (conn *SafeConn) releaseGlobalLock() {
	atomic.StoreInt32(&conn.globalLockStatus, 0)
}

func (conn *SafeConn) acquireReadLock() bool {
	if atomic.LoadInt32(&conn.globalLockStatus) == 1 {
		return false
	}
	atomic.AddInt32(&conn.readLockCount, 1)
	return true
}

func (conn *SafeConn) releaseReadLock() {
	atomic.AddInt32(&conn.readLockCount, -1)
}

// Starts a new transaction based on the current context. Will cancel if
// the context is closed/cancelled/done
func (conn *SafeConn) Begin(ctx context.Context) (*SafeTX, error) {
	lockAcquired := make(chan struct{})
	lockCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-lockCtx.Done():
			return
		default:
			if conn.acquireReadLock() {
				close(lockAcquired) // Lock acquired
			}
		}
	}()

	select {
	case <-lockAcquired:
		tx, err := conn.db.BeginTx(ctx, nil)
		if err != nil {
			conn.releaseReadLock()
			return nil, err
		}
		return &SafeTX{tx: tx, sc: conn}, nil
	case <-ctx.Done():
		cancel()
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

	stx.sc.releaseReadLock()
	return err
}

// Rollback aborts the transaction.
func (stx *SafeTX) Rollback() error {
	if stx.tx == nil {
		return errors.New("Cannot rollback without a transaction")
	}
	err := stx.tx.Rollback()
	stx.tx = nil
	stx.sc.releaseReadLock()
	return err
}

// Pause blocks new transactions for a backup.
func (conn *SafeConn) Pause() {
	for !conn.acquireGlobalLock() {
		// TODO: add a timeout?
	}
	fmt.Println("Global database lock acquired")
}

// Resume allows transactions to proceed.
func (conn *SafeConn) Resume() {
	conn.releaseGlobalLock()
	fmt.Println("Global database lock released")
}

// Close the database connection
func (conn *SafeConn) Close() error {
	conn.acquireGlobalLock()
	defer conn.releaseGlobalLock()
	return conn.db.Close()
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
