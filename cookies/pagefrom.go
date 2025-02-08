package cookies

import (
	"net/http"
	"net/url"
	"os"
	"time"
)

func CheckPageFrom(w http.ResponseWriter, r *http.Request) string {
	pageFromCookie, err := r.Cookie("pagefrom")
	if err != nil {
		return "/"
	}
	pageFrom := pageFromCookie.Value
	deleteCookie := &http.Cookie{Name: "pagefrom", Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, deleteCookie)
	return pageFrom
}

func SetPageFrom(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return
	}
	var pageFrom string
	expectedHost := os.Getenv("TRUSTED_HOST")
	if parsedURL.Path == "" || parsedURL.Host != expectedHost {
		pageFrom = "/"
	} else {
		pageFrom = parsedURL.Path
	}
	pageFromCookie := &http.Cookie{Name: "pagefrom", Value: pageFrom}
	http.SetCookie(w, pageFromCookie)
}
