package web

import (
	"log"
	"net/http"

	"github.com/emredir/songme/common/auth"
	"github.com/emredir/songme/models"
)

// Middleware provides functions for authorizing users.
type Middleware struct {
	userStore models.UserStore
}

// Authorize ensures user logged in.
func (mw *Middleware) Authorize(next http.Handler) http.Handler {
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

		_, err = mw.userStore.ByUUID(uuid)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// PanicRecovery recovers panic.
func (mw *Middleware) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[PanicRecovery]:", err)
				http.Error(w, "Opps! Things went wrong. 500: Server error.", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
