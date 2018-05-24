package controllers

import (
	"net/http"

	"github.com/emredir/songme/common/auth"
)

// LogoutController logouts user.
func LogoutController() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.ClearSessionCookie(w)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}
