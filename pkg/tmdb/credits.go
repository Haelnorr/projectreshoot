package tmdb

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Credits struct {
	ID   int32  `json:"id"`
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}

type Cast struct {
	Adult        bool   `json:"adult"`
	Gender       int    `json:"gender"`
	ID           int32  `json:"id"`
	KnownFor     string `json:"known_for_department"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Popularity   int    `json:"popularity"`
	Profile      string `json:"profile_path"`
	CastID       int32  `json:"cast_id"`
	Character    string `json:"character"`
	CreditID     string `json:"credit_id"`
	Order        int    `json:"order"`
}

type Crew struct {
	Adult        bool   `json:"adult"`
	Gender       int    `json:"gender"`
	ID           int32  `json:"id"`
	KnownFor     string `json:"known_for_department"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Popularity   int    `json:"popularity"`
	Profile      string `json:"profile_path"`
	CreditID     string `json:"credit_id"`
	Department   string `json:"department"`
	Job          string `json:"job"`
}

func GetCredits(movieid int32, token string) (*Credits, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v/credits?language=en-US", movieid)
	data, err := tmdbGet(url, token)
	if err != nil {
		return nil, errors.Wrap(err, "tmdbGet")
	}
	credits := Credits{}
	json.Unmarshal(data, &credits)
	return &credits, nil
}
