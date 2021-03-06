package models

import "time"

// Song represents the song that we recommend in our app.
type Song struct {
	ID          string
	Title       string
	Artist      string
	SongURL     string
	ImageURL    *string // nullable in db.
	Description *string // nullable in db.
	Confirmed   bool
	CreatedAt   time.Time
	ConfirmedAt *time.Time // nullable in db.
	UserID      *string    // nullable in db.
	User        *User      // nullable in db.
}

// DateString represents created at in readable format.
func (s *Song) DateString() string {
	return s.CreatedAt.Format("02.01.2006")
}

// SongStore defines the interface used to interact with the songs datastore.
type SongStore interface {
	Create(song *Song) error
	Confirm(id string) error

	GetRandom() (*Song, error)
	All(confirmed bool, limit, offset int) ([]*Song, int, error)
	UserSongs(id string, limit, offset int) ([]*Song, int, error)

	Delete(id string) error
}
