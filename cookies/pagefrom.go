package cookies

import (
	"net/http"
	"net/url"
)

// Check the value of "pagefrom" cookie, delete the cookie, and return the value
func CheckPageFrom(w http.ResponseWriter, r *http.Request) string {
	pageFromCookie, err := r.Cookie("pagefrom")
	if err != nil {
		return "/"
	}
	pageFrom := pageFromCookie.Value
	DeleteCookie(w, pageFromCookie.Name, pageFromCookie.Path)
	return pageFrom
}

// Check the referer of the request, and if it matches the trustedHost, set
// the "pagefrom" cookie as the Path of the referer
func SetPageFrom(w http.ResponseWriter, r *http.Request, trustedHost string) {
	referer := r.Referer()
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return
	}
	// NOTE: its possible this could cause an infinite redirect
	// if that happens, will need to add a way to 'blacklist' certain paths
	// from being set here
	var pageFrom string
	if parsedURL.Path == "" || parsedURL.Host != trustedHost {
		pageFrom = "/"
	} else {
		pageFrom = parsedURL.Path
	}
	pageFromCookie := &http.Cookie{Name: "pagefrom", Value: pageFrom, Path: "/"}
	http.SetCookie(w, pageFromCookie)
}
