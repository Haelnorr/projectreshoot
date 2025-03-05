package tmdb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Result struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ResultMovies struct {
	Result
	Results []ResultMovie `json:"results"`
}
type ResultMovie struct {
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
}

func (movie *ResultMovie) GetPoster(image *Image, size string) string {
	base, err := url.Parse(image.SecureBaseURL)
	if err != nil {
		return ""
	}
	fullPath := path.Join(base.Path, size, movie.PosterPath)
	base.Path = fullPath
	return base.String()
}

func (movie *ResultMovie) ReleaseYear() string {
	if movie.ReleaseDate == "" {
		return ""
	} else {
		return "(" + movie.ReleaseDate[:4] + ")"
	}
}

// TODO: genres list https://developer.themoviedb.org/reference/genre-movie-list
// func (movie *ResultMovie) FGenres() string {
// 	genres := ""
// 	for _, genre := range movie.Genres {
// 		genres += genre.Name + ", "
// 	}
// 	return genres[:len(genres)-2]
// }

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
