package tmdb

import (
	"fmt"
	"net/url"
	"path"
)

func (movie *Movie) FRuntime() string {
	hours := movie.Runtime / 60
	mins := movie.Runtime % 60
	return fmt.Sprintf("%dh %02dm", hours, mins)
}

func (movie *Movie) GetPoster(image *Image, size string) string {
	base, err := url.Parse(image.SecureBaseURL)
	if err != nil {
		return ""
	}
	fullPath := path.Join(base.Path, size, movie.Poster)
	base.Path = fullPath
	return base.String()
}

func (movie *Movie) ReleaseYear() string {
	if movie.ReleaseDate == "" {
		return ""
	} else {
		return "(" + movie.ReleaseDate[:4] + ")"
	}
}

func (movie *Movie) FGenres() string {
	genres := ""
	for _, genre := range movie.Genres {
		genres += genre.Name + ", "
	}
	if len(genres) > 2 {
		return genres[:len(genres)-2]
	}
	return genres
}
