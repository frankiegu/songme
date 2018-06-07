package models

import (
	"crypto/md5"
	"fmt"
)

// User represents the user in our system.
type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Role         Role
}

// IsAdmin tells whether the user is admin.
func (u *User) IsAdmin() bool {
	return u.Role.Has(PermAdmin)
}

// AvatarURL returns user's gravatar URL
func (u *User) AvatarURL() string {
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", md5.Sum([]byte(u.Email)))
}

// UserStore defines the interface used to interact with the users datastore.
type UserStore interface {
	Create(email, username, passwordHash string, roleID int) (string, error)

	ByID(id string) (*User, error)
	ByEmail(email string) (*User, error)
	ByUsername(username string) (*User, error)
	All(limit, offset int) ([]*User, int, error)

	UpdatePassword(email, passwordHash string) error
	UpdateEmail(old, new string) error
}
