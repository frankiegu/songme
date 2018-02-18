package auth

import (
	"log"
	"net/http"

	"github.com/emre-demir/songme/models"
	"github.com/gorilla/securecookie"
)

var (
	cookieName string

	hashKey  []byte
	blockKey []byte
	s        *securecookie.SecureCookie
)

// SetSessionCookie sets session cookie.
func SetSessionCookie(w http.ResponseWriter, u *models.User) error {
	value := map[string]string{
		"uuid": u.UUID,
	}
	if encoded, err := s.Encode(cookieName, &value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println("[SetSessionCookie:]", err)
		return err
	}
	return nil
}

// GetSessionCookie decodes session cookie.
func GetSessionCookie(r *http.Request) (map[string]string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println("[GetSessionCookie:]", err)
		return nil, err
	}

	value := make(map[string]string)
	err = s.Decode(cookieName, cookie.Value, &value)
	if err != nil {
		log.Println("[GetSessionCookie:]", err)
		return nil, err
	}

	return value, nil
}

// ClearSessionCookie clears session cookie for given http.ResponseWriter
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func init() {
	cookieName = "SongmeSession"

	hashKey = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	s = securecookie.New(hashKey, blockKey)
}
