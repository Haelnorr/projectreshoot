package tmdb

import (
// "encoding/json"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int    `json:"id"`
	Logo          string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	ISO_3166_1 string `json:"iso_3166_1"`
	Name       string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	ISO_639_1   string `json:"iso_639_1"`
	Name        string `json:"name"`
}
