package psql

import (
	"database/sql"

	"github.com/emredir/songme/models"
)

// SongStore is a PostgreSQL specific implementation of the song datastore.
type SongStore struct {
	DB *sql.DB
}

var _ models.SongStore = &SongStore{}

/*
	Implementation of the models.SongStore Interface
*/

// Create inserts a new song into database.
func (s *SongStore) Create(song *models.Song) error {
	stmt := `
	INSERT INTO candidate_song (
		title, 
		author, 
		song_url,
		image_url, 
		description
	) 
	VALUES ($1, $2, $3, $4, $5);`
	_, err := s.DB.Exec(stmt, song.Title, song.Author, song.SongURL, song.ImageURL, song.Description)
	return err
}

// Confirm transfers a song from table candidate_song to production_song.
func (s *SongStore) Confirm(id string) error {
	stmt := `
	INSERT INTO production_song (
		title, 
		author, 
		song_url,
		image_url, 
		description
	) 
	SELECT title, author, song_url, image_url, description
	FROM candidate_song
	WHERE id = $1;`
	_, err := s.DB.Exec(stmt, id)
	if err == nil {
		s.DeleteCandidate(id)
	}
	return err
}

// GetRandom returns a random song from the database.
func (s *SongStore) GetRandom() (*models.Song, error) {
	song := models.Song{}
	row := s.DB.QueryRow("SELECT * FROM production_song ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(
		&song.ID,
		&song.Title,
		&song.Author,
		&song.SongURL,
		&song.ImageURL,
		&song.Description,
		&song.Recommended,
		&song.CreatedAt,
		&song.RecommendedAt,
	)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// Candidates returns all candidate songs.
func (s *SongStore) Candidates() ([]*models.Song, error) {
	table := "candidate_song"

	rows, err := s.DB.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []*models.Song{}
	for rows.Next() {
		song := models.Song{}
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.Author,
			&song.SongURL,
			&song.ImageURL,
			&song.Description,
			&song.Recommended,
			&song.CreatedAt,
			&song.RecommendedAt,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

// Productions returns all songs in production.
func (s *SongStore) Productions() ([]*models.Song, error) {
	table := "production_song"

	rows, err := s.DB.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []*models.Song{}
	for rows.Next() {
		song := models.Song{}
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.Author,
			&song.SongURL,
			&song.ImageURL,
			&song.Description,
			&song.Recommended,
			&song.CreatedAt,
			&song.RecommendedAt,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

// DeleteCandidate deletes song from table 'candidate_song'.
func (s *SongStore) DeleteCandidate(id string) error {
	stmt := "DELETE FROM candidate_song WHERE id = $1;"
	_, err := s.DB.Exec(stmt, id)
	return err
}

// DeleteProduction deletes song from table 'production_song.
func (s *SongStore) DeleteProduction(id string) error {
	stmt := "DELETE FROM production_song WHERE id = $1;"
	_, err := s.DB.Exec(stmt, id)
	return err
}
