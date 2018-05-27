package models

import "time"

// Song represents the song that we recommend in our app.
type Song struct {
	ID            string
	Title         string
	Author        string
	SongURL       string
	ImageURL      string
	Description   string
	Recommended   bool
	CreatedAt     time.Time
	RecommendedAt time.Time
}

// SongStore defines the interface used to interact with the songs datastore.
type SongStore interface {
	Create(song *Song) error
	Confirm(id string) error

	GetRandom() (*Song, error)
	Candidates() ([]*Song, error)
	Productions() ([]*Song, error)

	DeleteCandidate(id string) error
	DeleteProduction(id string) error
}
