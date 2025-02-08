package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"projectreshoot/cookies"
	"projectreshoot/view/component/form"
	"projectreshoot/view/page"
)

// TODO: here for testing only, move to database
type User struct {
	id       int
	username string
	password string
}

// TODO: here for testing only, move to database
func testUser() User {
	return User{id: 1, username: "Haelnorr", password: "test"}
}

func validateLogin(r *http.Request) (int, error) {
	formUsername := r.FormValue("username")
	formPassword := r.FormValue("password")
	// TODO: search database for username
	validUser := testUser()
	// TODO: check password is valid
	if formUsername != validUser.username || formPassword != validUser.password {
		return 0, errors.New("Username or password incorrect")
	}
	// TODO: return the users ID
	return validUser.id, nil
}

func checkRememberMe(r *http.Request) bool {
	rememberMe := r.FormValue("remember-me")
	if rememberMe == "on" {
		return true
	} else {
		return false
	}
}

func HandleLoginRequest() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			userID, err := validateLogin(r)
			if err != nil {
				// TODO: add debug log
				fmt.Printf("Login failed: %s\n", err)
				form.LoginForm(err.Error()).Render(r.Context(), w)
				return
			}

			// TODO: login success, use the userID to set the session
			rememberMe := checkRememberMe(r)
			fmt.Printf("Login success, user ID: %v - remember me?: %t\n", userID, rememberMe)

			pageFrom := cookies.CheckPageFrom(w, r)
			w.Header().Set("HX-Redirect", pageFrom)
		},
	)
}

func HandleLoginPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookies.SetPageFrom(w, r)
			page.Login().Render(r.Context(), w)
		},
	)
}
