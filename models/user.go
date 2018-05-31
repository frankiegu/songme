package models

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

// UserStore defines the interface used to interact with the users datastore.
type UserStore interface {
	Create(email, username, passwordHash string, roleID int) (string, error)

	ByID(id string) (*User, error)
	ByEmail(email string) (*User, error)
	ByUsername(username string) (*User, error)
}
