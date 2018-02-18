package models

// User holds information about our users.
type User struct {
	UUID         string
	Username     string
	Email        string
	PasswordHash string
	Role         string
}
