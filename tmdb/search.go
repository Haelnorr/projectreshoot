package tmdb

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type Result struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ResultMovies struct {
	Result
	Results []struct {
		Adult            bool   `json:"adult"`
		BackdropPath     string `json:"backdrop_path"`
		GenreIDs         []int  `json:"genre_ids"`
		ID               int32  `json:"id"`
		OriginalLanguage string `json:"original_language"`
		OriginalTitle    string `json:"original_title"`
		Overview         string `json:"overview"`
		Popularity       int    `json:"popularity"`
		PosterPath       string `json:"poster_path"`
		ReleaseDate      string `json:"release_date"`
		Title            string `json:"title"`
		Video            bool   `json:"video"`
		VoteAverage      int    `json:"vote_average"`
		VoteCount        int    `json:"vote_count"`
	} `json:"results"`
}

func SearchMovies(token string, query string, adult bool, page int) (*ResultMovies, error) {
	url := "https://api.themoviedb.org/3/search/movie" +
		fmt.Sprintf("?query=%s", url.QueryEscape(query)) +
		fmt.Sprintf("&include_adult=%t", adult) +
		fmt.Sprintf("&page=%v", page) +
		"&language=en-US"
	response, err := tmdbGet(url, token)
	if err != nil {
		return nil, errors.Wrap(err, "tmdbGet")
	}
	var results ResultMovies
	json.Unmarshal(response, &results)
	return &results, nil
}
