package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	_ "modernc.org/sqlite"
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
func SetupTestDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
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

	_, err = conn.Exec(schemaSQL)
	if err != nil {
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

	_, err = conn.Exec(dataSQL)
	if err != nil {
		return nil, errors.Wrap(err, "tx.Exec")
	}
	return conn, nil
}
