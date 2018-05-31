package psql

import (
	"database/sql"

	"github.com/emredir/songme/models"
)

// UserStore is a PostgreSQL specific implementation of the user datastore.
type UserStore struct {
	DB *sql.DB
}

var _ models.UserStore = &UserStore{}

/*
	Implementation of the models.UserStore Interface
*/

// Create inserts a new user into database.
func (s *UserStore) Create(email, username, passwordHash string, roleID int) (string, error) {
	stmt := `
	INSERT INTO users (
		email, 
		username, 
		password_hash,
		role_id
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id string
	err := s.DB.QueryRow(stmt, email, username, passwordHash, roleID).Scan(&id)
	return id, err
}

// ByID returns user with given id.
func (s *UserStore) ByID(id string) (*models.User, error) {
	user := models.User{}
	row := s.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.password_hash, r.id, r.name, r.default_role, r.permissions 
		FROM users AS u 
		INNER JOIN roles AS r 
		ON u.role_id = r.id 
		WHERE u.id = $1`, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Default,
		&user.Role.Permissions,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail returns user with given email.
func (s *UserStore) ByEmail(email string) (*models.User, error) {
	user := models.User{}
	row := s.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.password_hash, r.id, r.name, r.default_role, r.permissions 
		FROM users AS u 
		INNER JOIN roles AS r 
		ON u.role_id = r.id 
		WHERE u.email = $1`, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Default,
		&user.Role.Permissions,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByUsername returns user with given username (both username and email are accepted).
func (s *UserStore) ByUsername(username string) (*models.User, error) {
	user := models.User{}
	row := s.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.password_hash, r.id, r.name, r.default_role, r.permissions 
		FROM users AS u 
		INNER JOIN roles AS r 
		ON u.role_id = r.id 
		WHERE u.username = $1 OR u.email = $1`, username)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Default,
		&user.Role.Permissions,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
