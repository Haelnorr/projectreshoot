package tmdb

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func tmdbGet(url string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}
	return body, nil
}
