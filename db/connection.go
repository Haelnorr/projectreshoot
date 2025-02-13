package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// Returns a database connection handle for the Turso DB
func ConnectToDatabase(dbName string) (*sql.DB, error) {
	file := fmt.Sprintf("file:%s.db", dbName)
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}
	return db, nil
}
