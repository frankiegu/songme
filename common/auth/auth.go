package auth

import (
	"errors"
	"log"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/common/utility"
	"github.com/emre-demir/songme/models"
)

// LoginUserWith checks the given user credentials.
func LoginUserWith(username string, password string, ev *env.Vars) (*models.User, error) {
	u, err := ev.DB.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if utility.SHA256String(password) == u.PasswordHash {
		log.Println("User logged in:", u.Username)
		return u, nil
	}

	log.Println("Wrong credentials for user:", username)
	return nil, errors.New("Wrong credentials")
}
