package handler

import (
	"net/http"
	"projectreshoot/config"
	"projectreshoot/tmdb"
	"projectreshoot/view/component/search"
	"projectreshoot/view/page"

	"github.com/rs/zerolog"
)

func SearchMovies(
	logger *zerolog.Logger,
	config *config.Config,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			query := r.FormValue("search")
			if query == "" {
				w.WriteHeader(http.StatusOK)
				return
			}
			movies, err := tmdb.SearchMovies(config.TMDBToken, query, false, 1)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			search.MovieResults(movies, &config.TMDBConfig.Image).Render(r.Context(), w)
		},
	)
}

func MoviesPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			page.Movies().Render(r.Context(), w)
		},
	)
}
