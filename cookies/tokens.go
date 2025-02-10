package cookies

import (
	"net/http"
	"projectreshoot/config"
	"time"
)

// Get the value of the access and refresh tokens
func GetTokens(
	w http.ResponseWriter,
	r *http.Request,
) (acc string, ref string) {
	accCookie, accErr := r.Cookie("access")
	refCookie, refErr := r.Cookie("refresh")
	var (
		accStr string = ""
		refStr string = ""
	)
	if accErr == nil {
		accStr = accCookie.Value
	}
	if refErr == nil {
		refStr = refCookie.Value
	}
	return accStr, refStr
}

// Set a token with the provided details
func SetToken(
	w http.ResponseWriter,
	r *http.Request,
	config *config.Config,
	token string,
	scope string,
	exp int64,
) {
	tokenCookie := &http.Cookie{
		Name:     scope,
		Value:    token,
		Path:     "/",
		Expires:  time.Unix(exp, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   config.SSL,
	}
	http.SetCookie(w, tokenCookie)
}
