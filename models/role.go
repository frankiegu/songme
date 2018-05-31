package models

// Permission represents different type of permissions in our system.
type Permission int

const (
	// PermUser is for user permission
	PermUser Permission = 1

	// PermModerate is for moderate permission
	PermModerate Permission = 2

	// PermAdmin is for admin permission
	PermAdmin Permission = 4
)

// Role represents the user roles in our system.
type Role struct {
	ID          int
	Name        string
	Default     bool
	Permissions Permission
}

// Has checks if the role has given permission.
func (r Role) Has(p Permission) bool {
	return (r.Permissions & p) == p
}

// RoleStore defines the interface used to interact with the roles datastore.
type RoleStore interface {
	ByName(name string) (*Role, error)
}
