package psql

import (
	"database/sql"

	"github.com/emredir/songme/models"
)

// RoleStore is a PostgreSQL specific implementation of the role datastore.
type RoleStore struct {
	DB *sql.DB
}

var _ models.RoleStore = &RoleStore{}

/*
	Implementation of the models.RoleStore Interface
*/

// ByName returns user with given id.
func (s *RoleStore) ByName(name string) (*models.Role, error) {
	role := models.Role{}
	row := s.DB.QueryRow("SELECT * FROM roles WHERE name = $1", name)
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Default,
		&role.Permissions,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
