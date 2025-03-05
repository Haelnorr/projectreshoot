package tmdb

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Config struct {
	Image Image `json:"images"`
}

type Image struct {
	BaseURL       string   `json:"base_url"`
	SecureBaseURL string   `json:"secure_base_url"`
	BackdropSizes []string `json:"backdrop_sizes"`
	LogoSizes     []string `json:"logo_sizes"`
	PosterSizes   []string `json:"poster_sizes"`
	ProfileSizes  []string `json:"profile_sizes"`
	StillSizes    []string `json:"still_sizes"`
}

func GetConfig(token string) (*Config, error) {
	url := "https://api.themoviedb.org/3/configuration"
	data, err := tmdbGet(url, token)
	if err != nil {
		return nil, errors.Wrap(err, "tmdbGet")
	}
	config := Config{}
	json.Unmarshal(data, &config)
	return &config, nil
}
