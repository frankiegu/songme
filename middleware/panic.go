package middleware

import (
	"log"
	"net/http"
)

// PanicRecovery recovers panic occurs in next handlers.
func PanicRecovery(next http.Handler) http.Handler {
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
