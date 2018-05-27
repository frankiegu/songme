package models

// User represents the user in our system.
type User struct {
	UUID         string
	Username     string
	Email        string
	PasswordHash string
	Role         string
}

// UserStore defines the interface used to interact with the users datastore.
type UserStore interface {
	ByUUID(uuid string) (*User, error)
	ByUsername(username string) (*User, error)
}
