package interactor

import (
	"errors"

	"github.com/emredir/songme/models"
)

var (
	// ErrInternalSong should be returned when an internal error occurs.
	ErrInternalSong = errors.New("Error: Song interactor")

	// ErrSongExists should be returned when a song already exists.
	ErrSongExists = errors.New("Sorry, that song has already been in our system. Try another?")

	// ErrNoSongs should be returned when the there are no songs in the result.
	ErrNoSongs = errors.New("Sorry, it seems we can't find what you're looking for")
)

// Song performs database operations for usecases.
type Song struct {
	models.SongStore
}

// Create creates a new song.
func (s *Song) Create(song *models.Song) error {
	err := s.SongStore.Create(song)

	/*
		Can't use this because of migrate package!

		// Check for the type of postgresql error.
		if pqError, ok := err.(*pq.Error); ok {
			// Error '23505' is unique violation
			if pqError.Code == "23505" {
				return ErrSongExists
			}
		}
	*/

	// TODO: make search about it.
	if err != nil {
		return ErrSongExists
	}
	return nil
}

// All returns all songs by considering their confirmation status.
func (s *Song) All(confirmed bool, limit, offset int) ([]*models.Song, int, error) {
	songs, totalCount, err := s.SongStore.All(confirmed, limit, offset)
	if err != nil {
		return nil, totalCount, err
	}
	if len(songs) < 1 {
		return nil, totalCount, ErrNoSongs
	}
	return songs, totalCount, nil
}
