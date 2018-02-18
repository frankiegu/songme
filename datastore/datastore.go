package datastore

import (
	"errors"

	"github.com/emre-demir/songme/models"
)

// Enumerator for supported databases.
const (
	MySQL uint8 = iota
	PostgreSQL
)

// Datastore defines necessary functions for managing databases.
type Datastore interface {
	Close()
	Subscribe(*models.Subscriber) error
	CreateSong(*models.Song) error
	GetRandomSong() *models.Song
	GetSongsFrom(string) ([]*models.Song, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserByUUID(string) (*models.User, error)
}

// NewDatastore returns new datastore instance.
func NewDatastore(dbtype uint8, config Config) (Datastore, error) {
	switch dbtype {
	case PostgreSQL:
		return newPQDatastore(config)
	}
	return nil, errors.New("Unsupported database")
}
