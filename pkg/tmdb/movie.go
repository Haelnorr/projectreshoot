package tmdb

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Movie struct {
	Adult               bool                `json:"adult"`
	Backdrop            string              `json:"backdrop_path"`
	Collection          string              `json:"belongs_to_collection"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	ID                  int32               `json:"id"`
	IMDbID              string              `json:"imdb_id"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float32             `json:"popularity"`
	Poster              string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
}

func GetMovie(id int32, token string) (*Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v?language=en-US", id)
	data, err := tmdbGet(url, token)
	if err != nil {
		return nil, errors.Wrap(err, "tmdbGet")
	}
	movie := Movie{}
	json.Unmarshal(data, &movie)
	return &movie, nil
}
