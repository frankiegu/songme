package interactor

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/emredir/songme"
	"github.com/emredir/songme/models"
)

var (
	// ErrInternal should be returned when an internal error occurs.
	ErrInternal = errors.New("Error: Authentication interactor")

	// ErrEmailExists should be returned when an email already exists.
	ErrEmailExists = errors.New("Sorry, that email has already been taken")

	// ErrUsernameExists should be returned when a username already exists.
	ErrUsernameExists = errors.New("Sorry, that username's taken. Try another?")

	// ErrWrongCredentials should be returned when a user enters wrong credentials.
	ErrWrongCredentials = errors.New("Your credentials are incorrect")
)

// Auth interacts with database for operations such as sign up, sign in etc.
type Auth struct {
	UserStore models.UserStore
	RoleStore models.RoleStore
}

// Signup sign up a new user to the system.
func (a *Auth) Signup(email, username, password string) (*models.User, error) {
	user, _ := a.UserStore.ByUsername(username)
	if user != nil {
		return nil, ErrUsernameExists
	}

	user, _ = a.UserStore.ByEmail(email)
	if user != nil {
		return nil, ErrEmailExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), songme.GetConfig().HashCost)
	if err != nil {
		return nil, ErrInternal
	}

	var role *models.Role
	if email == songme.GetConfig().SongmeAdmin {
		role, err = a.RoleStore.ByName("Administrator")
	} else {
		role, err = a.RoleStore.ByName("User")
	}
	if err != nil {
		return nil, ErrInternal
	}

	id, err := a.UserStore.Create(email, username, string(passwordHash), role.ID)
	if err != nil {
		return nil, ErrInternal
	}

	log.Println("[interactor.Signup]: User registered", username)

	return &models.User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: string(passwordHash),
		Role:         *role,
	}, nil
}

// Signin sign in a user to the system.
func (a *Auth) Signin(username, password string) (*models.User, error) {
	user, err := a.UserStore.ByUsername(username)
	if err != nil {
		return nil, ErrWrongCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrWrongCredentials
	}

	log.Println("[interactor.Signin]: User logged in", username)

	return user, nil
}
