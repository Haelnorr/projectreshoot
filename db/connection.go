package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	_ "modernc.org/sqlite"
)

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
