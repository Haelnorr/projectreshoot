package tmdb

import (
	"fmt"
	"net/url"
	"path"
)

func FormatRuntime(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%dh %02dm", hours, mins)
}

func GetPoster(image *Image, size, imgpath string) string {
	base, err := url.Parse(image.SecureBaseURL)
	if err != nil {
		return ""
	}
	fullPath := path.Join(base.Path, size, imgpath)
	base.Path = fullPath
	return base.String()
}
