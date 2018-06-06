package web

import (
	"log"
	"net/http"

	"github.com/emredir/songme/internal/context"
	"github.com/emredir/songme/internal/cookie"
	"github.com/emredir/songme/models"
)

// Middleware provides functions for authorizing users.
type Middleware struct {
	userStore models.UserStore
}

// UserViaSession retrieves the current user by the session cookie
// and set it in the request context.
func (mw *Middleware) UserViaSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := cookie.Get(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		id, ok := cookieValue["id"]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		user, err := mw.userStore.ByID(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithUser(r.Context(), user))
		next.ServeHTTP(w, r)
	})
}

// Authorize ensures user logged in.
func (mw *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Admin ensures user has admin privileges.
func (mw *Middleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if !user.IsAdmin() {
			http.Redirect(w, r, "/", http.StatusFound)
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
				log.Println("[web.PanicRecovery]:", err)
				http.Error(w, "Opps! Things went wrong. 500: Server error.", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
