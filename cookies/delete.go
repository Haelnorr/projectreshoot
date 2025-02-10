package cookies

import (
	"net/http"
	"time"
)

// Tell the browser to delete the cookie matching the name provided
// Path must match the original set cookie for it to delete
func DeleteCookie(w http.ResponseWriter, name string, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    path,
		Expires: time.Unix(0, 0), // Expire in the past
		MaxAge:  -1,              // Immediately expire
	})
}
