package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"projectreshoot/db"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

func findSQLFile(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, filename)); err == nil {
			return filepath.Join(dir, filename), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // Reached root
			return "", errors.New(fmt.Sprintf("Unable to locate %s", filename))
		}
		dir = parent
	}
}

// SetupTestDB initializes a test SQLite database with mock data
// Make sure to call DeleteTestDB when finished to cleanup
func SetupTestDB(ctx context.Context) (*db.SafeConn, error) {
	dbfile, err := sql.Open("sqlite3", "file:.projectreshoot-test-database.db")
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}
	conn := db.MakeSafe(dbfile)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Begin")
	}
	// Setup the test database
	schemaPath, err := findSQLFile("schema.sql")
	if err != nil {
		return nil, errors.Wrap(err, "findSchema")
	}

	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}
	schemaSQL := string(sqlBytes)

	_, err = tx.Exec(ctx, schemaSQL)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "tx.Exec")
	}
	// Load the test data
	dataPath, err := findSQLFile("testdata.sql")
	if err != nil {
		return nil, errors.Wrap(err, "findSchema")
	}
	sqlBytes, err = os.ReadFile(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}
	dataSQL := string(sqlBytes)

	_, err = tx.Exec(ctx, dataSQL)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "tx.Exec")
	}
	tx.Commit()
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
