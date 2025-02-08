package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"projectreshoot/cookies"
	"projectreshoot/db"
	"projectreshoot/view/component/form"
	"projectreshoot/view/page"

	"github.com/pkg/errors"
)

func validateLogin(conn *sql.DB, r *http.Request) (db.User, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")

	user, err := db.GetUserFromUsername(conn, formUsername)
	if err != nil {
		return db.User{}, errors.Wrap(err, "db.GetUserFromUsername")
	}

	err = user.CheckPassword(formPassword)
	if err != nil {
		return db.User{}, errors.New("Username or password incorrect")
	}
	return user, nil
}

func checkRememberMe(r *http.Request) bool {
	rememberMe := r.FormValue("remember-me")
	if rememberMe == "on" {
		return true
	} else {
		return false
	}
}

func HandleLoginRequest(conn *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			user, err := validateLogin(conn, r)
			if err != nil {
				form.LoginForm(err.Error()).Render(r.Context(), w)
				return
			}

			// TODO: login success, use the userID to set the session
			rememberMe := checkRememberMe(r)
			fmt.Printf(
				"Login success, user: %v - remember me?: %t\n",
				user.Username,
				rememberMe,
			)

			pageFrom := cookies.CheckPageFrom(w, r)
			w.Header().Set("HX-Redirect", pageFrom)
		},
	)
}

func HandleLoginPage(trustedHost string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookies.SetPageFrom(w, r, trustedHost)
			page.Login().Render(r.Context(), w)
		},
	)
}
