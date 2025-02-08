package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func ConnectToDatabase(primaryUrl *string, authToken *string) (*sql.DB, error) {
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", *primaryUrl, *authToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}
	return db, nil
}
