package middleware

import (
	"net/http"

	"github.com/emredir/songme/common/env"

	"github.com/emredir/songme/common/auth"
)

// Authorize ensures user logged in.
func Authorize(next http.Handler, ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := auth.GetSessionCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		uuid, ok := cookieValue["uuid"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, err = ev.DB.GetUserByUUID(uuid)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
