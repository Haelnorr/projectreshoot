package main

import (
	"database/sql"
	"net/http"

	"projectreshoot/server"
)

// This function will only be called if the --test commandline flag is set.
// After the function finishes the application will close.
// Running command `make test` will run the test using port 3232 to avoid
// conflicts on the default 3333. Useful for testing things out during dev
func test(config *server.Config, conn *sql.DB, srv *http.Server) {
}
