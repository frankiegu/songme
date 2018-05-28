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

// ByUsername returns user with given username.
func (s *UserStore) ByUsername(username string) (*models.User, error) {
	user := models.User{}
	row := s.DB.QueryRow("SELECT * FROM users WHERE username = $1 OR email = $1", username)
	err := row.Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByUUID returns user with given uuid.
func (s *UserStore) ByUUID(uuid string) (*models.User, error) {
	user := models.User{}
	row := s.DB.QueryRow("SELECT * FROM users WHERE uuid = $1", uuid)
	err := row.Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
