package main

import (
	"net/http"

	"projectreshoot/pkg/config"
	"projectreshoot/pkg/db"
	"projectreshoot/pkg/tmdb"

	"github.com/rs/zerolog"
)

// This function will only be called if the --test commandline flag is set.
// After the function finishes the application will close.
// Running command `make tester` will run the test using port 3232 to avoid
// conflicts on the default 3333. Useful for testing things out during dev.
// If you add code here, remember to run:
// `git update-index --assume-unchanged tester.go` to avoid tracking changes
func test(
	config *config.Config,
	logger *zerolog.Logger,
	conn *db.SafeConn,
	srv *http.Server,
) {
	query := "a few good men"
	search, err := tmdb.SearchMovies(config.TMDBToken, query, false, 1)
	if err != nil {
		logger.Error().Err(err).Msg("error occured")
		return
	}
	logger.Info().Interface("results", search).Msg("search results received")
}
