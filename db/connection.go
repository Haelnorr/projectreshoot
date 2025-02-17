package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	_ "modernc.org/sqlite"
)

type SafeConn struct {
	db               *sql.DB
	readLockCount    uint32
	globalLockStatus uint32
	logger           *zerolog.Logger
}

func MakeSafe(db *sql.DB, logger *zerolog.Logger) *SafeConn {
	return &SafeConn{db: db, logger: logger}
}

// Extends sql.Tx for use with SafeConn
type SafeTX struct {
	tx *sql.Tx
	sc *SafeConn
}

func (conn *SafeConn) acquireGlobalLock() bool {
	if conn.readLockCount > 0 || conn.globalLockStatus == 1 {
		return false
	}
	conn.globalLockStatus = 1
	conn.logger.Debug().Uint32("global_lock_status", conn.globalLockStatus).
		Msg("Global lock acquired")
	return true
}

func (conn *SafeConn) releaseGlobalLock() {
	conn.globalLockStatus = 0
	conn.logger.Debug().Uint32("global_lock_status", conn.globalLockStatus).
		Msg("Global lock released")
}

func (conn *SafeConn) acquireReadLock() bool {
	if conn.globalLockStatus == 1 {
		return false
	}
	conn.readLockCount += 1
	conn.logger.Debug().Uint32("read_lock_count", conn.readLockCount).
		Msg("Read lock acquired")
	return true
}

func (conn *SafeConn) releaseReadLock() {
	conn.readLockCount -= 1
	conn.logger.Debug().Uint32("read_lock_count", conn.readLockCount).
		Msg("Read lock released")
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
		// TODO: failed to acquire lock: print info with readLockCount
		// every second, or update it dynamically
	}
	// force logger to log to Stdout
	log := conn.logger.With().Logger().Output(os.Stdout)
	log.Info().Msg("Global database lock acquired")
}

// Resume allows transactions to proceed.
func (conn *SafeConn) Resume() {
	conn.releaseGlobalLock()
	// force logger to log to Stdout
	log := conn.logger.With().Logger().Output(os.Stdout)
	log.Info().Msg("Global database lock released")
}

// Close the database connection
func (conn *SafeConn) Close() error {
	conn.logger.Debug().Msg("Acquiring global lock for connection close")
	conn.acquireGlobalLock()
	defer conn.releaseGlobalLock()
	conn.logger.Debug().Msg("Closing database connection")
	return conn.db.Close()
}

// Returns a database connection handle for the DB
func ConnectToDatabase(dbName string, logger *zerolog.Logger) (*SafeConn, error) {
	file := fmt.Sprintf("file:%s.db", dbName)
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	conn := MakeSafe(db, logger)

	return conn, nil
}
