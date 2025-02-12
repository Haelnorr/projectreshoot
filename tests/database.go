package tests

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB initializes a test SQLite database with mock data
// Make sure to call DeleteTestDB when finished to cleanup
func SetupTestDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "file:.projectreshoot-test-database.db")
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}
	// Create the test database
	_, err = conn.Exec(`
CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        username TEXT NOT NULL,
        password_hash TEXT,
        created_at INTEGER DEFAULT (unixepoch())
);
INSERT INTO users VALUES(1,'testuser','hashedpassword',1738995274);

CREATE TABLE IF NOT EXISTS jwtblacklist (
    jti TEXT PRIMARY KEY CHECK(jti GLOB '[0-9a-fA-F-]*'),
    exp INTEGER NOT NULL
) STRICT;

	`)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Exec")
	}
	return conn, nil
}

// Deletes the test database from disk
func DeleteTestDB() error {
	fileName := ".projectreshoot-test-database.db"

	// Attempt to remove the file
	err := os.Remove(fileName)
	if err != nil {
		return errors.Wrap(err, "os.Remove")
	}

	return nil
}
