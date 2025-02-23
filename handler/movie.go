package handler

import (
	"net/http"
	"projectreshoot/config"
	"projectreshoot/tmdb"
	"projectreshoot/view/page"
	"strconv"

	"github.com/rs/zerolog"
)

func Movie(
	logger *zerolog.Logger,
	config *config.Config,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("movie_id")
			movie_id, err := strconv.ParseInt(id, 10, 32)
			if err != nil {
				ErrorPage(http.StatusNotFound, w, r)
				logger.Error().Err(err).Str("movie_id", id).
					Msg("Error occured getting the movie")
				return
			}
			movie, err := tmdb.GetMovie(int32(movie_id), config.TMDBToken)
			if err != nil {
				ErrorPage(http.StatusInternalServerError, w, r)
				logger.Error().Err(err).Int32("movie_id", int32(movie_id)).
					Msg("Error occured getting the movie")
				return
			}
			page.Movie(movie, &config.TMDBConfig.Image).Render(r.Context(), w)
		},
	)
}
