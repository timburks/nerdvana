package user

import (
	"net/http"
)

// User represents a user of the application.
type User struct {
	Email      string
	Nickname   string
	AuthDomain string
	ID         string
	Admin      bool
}

// Current returns the currently logged-in user,
// or nil if the user is not signed in.
func Current(r *http.Request) *User {
	h := r.Header
	u := &User{
		Email:      h.Get("X-Appengine-User-Email"),
		Nickname:   h.Get("X-Appengine-User-Nickname"),
		AuthDomain: h.Get("X-Appengine-Auth-Domain"),
		ID:         h.Get("X-Appengine-User-Id"),
		Admin:      h.Get("X-Appengine-User-Is-Admin") == "1",
	}
	if u.Email == "" {
		return nil
	}
	return u
}

// LoginURL returns a URL that, when visited, prompts the user to sign in.
func LoginURL() string {
	return "/login"
}

// LogoutURL returns a URL that, when visited, signs the user out.
func LogoutURL() string {
	return "https://accounts.google.com/SignOutOptions?hl=en&continue=https://google.com"
}
