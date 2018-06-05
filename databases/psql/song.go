package psql

import (
	"database/sql"

	"github.com/emredir/songme/models"
)

var _ models.SongStore = &SongStore{}

// SongStore is a PostgreSQL specific implementation of the song datastore.
type SongStore struct {
	DB *sql.DB
}

// recommender is a helper function for fetching user who recommends the song.
func (s *SongStore) recommender(song *models.Song) error {
	if song.UserID == nil {
		return nil // This is not an error, just return from the function.
	}
	us := &UserStore{DB: s.DB}
	user, err := us.ByID(*song.UserID)
	song.User = user
	return err
}

/*
	Implementation of the models.SongStore Interface
*/

// Create inserts a new song into database.
func (s *SongStore) Create(song *models.Song) error {
	stmt := `
	INSERT INTO songs (
		title, 
		artist, 
		song_url,
		image_url, 
		description,
		user_id
	) 
	VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := s.DB.Exec(stmt, song.Title, song.Artist, song.SongURL, song.ImageURL, song.Description, song.UserID)
	return err
}

// Confirm sets confirmed true for the given song.
func (s *SongStore) Confirm(id string) error {
	stmt := "UPDATE songs SET confirmed = TRUE, confirmed_at = NOW() WHERE id = $1;"
	_, err := s.DB.Exec(stmt, id)
	return err
}

/*
	This will raise the sql error: "storing driver.Value type <nil> into type *string"
	while scanning if the user_id is null.

	`
	SELECT
	s.id, s.title, s.artist, s.song_url, s.image_url, s.description, s.confirmed, s.created_at, s.confirmed_at, s.user_id,
	u.email, u.username, u.password_hash, u.role_id,
	r.name, r.default_role, r.permissions
	FROM songs AS s
	LEFT JOIN users AS u
		ON s.user_id = u.id
	LEFT JOIN roles AS r
		ON u.role_id = r.id
	ORDER BY RANDOM()
	LIMIT 1;
	`
*/

// GetRandom returns a random song from the database.
func (s *SongStore) GetRandom() (*models.Song, error) {
	stmt := "SELECT * FROM songs WHERE confirmed = TRUE ORDER BY RANDOM() LIMIT 1;"
	row := s.DB.QueryRow(stmt)
	song := models.Song{}
	err := row.Scan(
		&song.ID,
		&song.Title,
		&song.Artist,
		&song.SongURL,
		&song.ImageURL,
		&song.Description,
		&song.Confirmed,
		&song.CreatedAt,
		&song.ConfirmedAt,
		&song.UserID,
	)
	if err != nil {
		return nil, err
	}
	err = s.recommender(&song)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// All returns all songs, and their total count by considering confirmation status.
func (s *SongStore) All(confirmed bool, limit, offset int) ([]*models.Song, int, error) {
	// As we know, OFFSET solution leads to poor performance.
	// Refactor if necessary.
	stmt := `
	SELECT *, COUNT(confirmed) OVER() AS total_count
	FROM songs 
	WHERE confirmed = $1
	ORDER BY created_at DESC 
	LIMIT $2
	OFFSET $3;`
	rows, err := s.DB.Query(stmt, confirmed, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	totalCount := 0
	songs := []*models.Song{}
	for rows.Next() {
		song := models.Song{}
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.Artist,
			&song.SongURL,
			&song.ImageURL,
			&song.Description,
			&song.Confirmed,
			&song.CreatedAt,
			&song.ConfirmedAt,
			&song.UserID,
			&totalCount,
		)
		if err != nil {
			return nil, 0, err
		}
		err = s.recommender(&song)
		if err != nil {
			return nil, 0, err
		}
		songs = append(songs, &song)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return songs, totalCount, nil
}

// Delete deletes song from the database.
func (s *SongStore) Delete(id string) error {
	stmt := "DELETE FROM songs WHERE id = $1;"
	_, err := s.DB.Exec(stmt, id)
	return err
}
