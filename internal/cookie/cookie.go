package cookie

import (
	"log"
	"net/http"

	"github.com/emredir/songme"

	"github.com/gorilla/securecookie"
)

var (
	cookieName string

	hashKey  string
	blockKey string
	s        *securecookie.SecureCookie
)

func init() {
	cookieName = songme.GetConfig().CookieName

	hashKey = songme.GetConfig().CookieHashKey
	if hashKey == "" {
		hashKey = string(securecookie.GenerateRandomKey(64))
	}

	blockKey = songme.GetConfig().CookieBlockKey
	if blockKey == "" {
		blockKey = string(securecookie.GenerateRandomKey(32))
	}

	s = securecookie.New([]byte(hashKey), []byte(blockKey))
}

// Set sets the session cookie.
func Set(w http.ResponseWriter, id string) error {
	value := map[string]string{
		"id": id,
	}
	if encoded, err := s.Encode(cookieName, &value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println("[cookie.Set]:", err)
		return err
	}
	return nil
}

// Get returns session cookie.
func Get(r *http.Request) (map[string]string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	value := make(map[string]string)

	err = s.Decode(cookieName, cookie.Value, &value)
	if err != nil {
		log.Println("[cookie.Get]:", err)
		return nil, err
	}

	return value, nil
}

// Clear clears the session cookie.
func Clear(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
