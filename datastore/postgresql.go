package datastore

import (
	"database/sql"
	"log"
	// Run postgresql driver.
	_ "github.com/lib/pq"

	"github.com/emredir/songme/models"
)

// pQDatastore handles CRUD operations on PostgreSQL database.
type pQDatastore struct {
	*sql.DB
}

// Ensure pQDatastore conforms Datastore interface.
var _ Datastore = &pQDatastore{}

// Close closes database.
func (pq *pQDatastore) Close() {
	pq.DB.Close()
}

// Subscribe inserts given subscriber into database.
func (pq *pQDatastore) Subscribe(subscriber *models.Subscriber) error {
	tx, err := pq.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO subscriber (name, email) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(subscriber.Name, subscriber.Email)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CreateSong inserts given song into database.
func (pq *pQDatastore) CreateSong(song *models.Song) error {
	stmt := `
	INSERT INTO candidate_song (
		title, 
		author, 
		song_url,
		image_url, 
		description
	) 
	VALUES ($1, $2, $3, $4, $5);`
	_, err := pq.Exec(stmt, song.Title, song.Author, song.SongURL, song.ImageURL, song.Description)
	return err
}

// SelectSong transfers song with given id from
// table 'candidate_song' to 'production_song'.
func (pq *pQDatastore) SelectSong(id string) error {
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
	_, err := pq.Exec(stmt, id)
	if err == nil {
		pq.DeleteCandidateSong(id)
	}

	return err
}

// DeleteSong deletes song from table 'production_song'.
func (pq *pQDatastore) DeleteProductionSong(id string) error {
	stmt := "DELETE FROM production_song WHERE id = $1;"
	_, err := pq.Exec(stmt, id)
	return err
}

// DeleteCandidateSong deletes song from table 'candidate_song'.
func (pq *pQDatastore) DeleteCandidateSong(id string) error {
	stmt := "DELETE FROM candidate_song WHERE id = $1;"
	_, err := pq.Exec(stmt, id)
	return err
}

// GetRandomSong returns a randomly selected song.
func (pq *pQDatastore) GetRandomSong() *models.Song {
	song := models.Song{}

	row := pq.QueryRow("SELECT * FROM production_song ORDER BY RANDOM() LIMIT 1")
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
		log.Println("GetRandomSong:", err)
	}

	return &song
}

// GetSongsFrom returns all songs from given table.
func (pq *pQDatastore) GetSongsFrom(table string) ([]*models.Song, error) {
	rows, err := pq.Query("SELECT * FROM " + table)
	if err != nil {
		log.Println("[GetSongsFrom]:", table, err)
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
			log.Println("[GetSongsFrom]:", table, err)
			return nil, err
		}
		songs = append(songs, &song)
	}

	err = rows.Err()
	if err != nil {
		log.Println("[GetSongsFrom]:", table, err)
		return nil, err
	}

	return songs, nil
}

// GetUserByUsername returns user match with given username.
func (pq *pQDatastore) GetUserByUsername(username string) (*models.User, error) {
	user := models.User{}

	row := pq.QueryRow("SELECT * FROM users WHERE username = $1 OR email = $1", username)
	err := row.Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		log.Println("GetUserByUsername:", err)
		return nil, err
	}

	return &user, nil
}

// GetUserByUUID returns user match with given uuid.
func (pq *pQDatastore) GetUserByUUID(uuid string) (*models.User, error) {
	user := models.User{}

	row := pq.QueryRow("SELECT * FROM users WHERE uuid = $1", uuid)
	err := row.Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		log.Println("GetUserByUUID:", err)
		return nil, err
	}

	return &user, nil
}

// newPQDatastore returns new PQDatastore instance.
func newPQDatastore(config Config) (*pQDatastore, error) {
	var prepareDBStatements = []string{
		`CREATE TABLE IF NOT EXISTS subscriber (
			id SERIAL,
			name VARCHAR(64) NOT NULL,
			email VARCHAR(64) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			UNIQUE (email),
			PRIMARY KEY (id)
		);`,
		`CREATE TABLE IF NOT EXISTS production_song (
			id SERIAL,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			song_url VARCHAR(255) NOT NULL,
			image_url VARCHAR(255),
			description VARCHAR(280),
			recommended BOOLEAN NOT NULL DEFAULT false,			
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			recommended_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			UNIQUE (title),
			UNIQUE (song_url),
			PRIMARY KEY (id)
		);`,
		`CREATE TABLE IF NOT EXISTS candidate_song (
			id SERIAL,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			song_url VARCHAR(255) NOT NULL,
			image_url VARCHAR(255),
			description VARCHAR(280),
			recommended BOOLEAN NOT NULL DEFAULT false,			
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			recommended_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			UNIQUE (title),
			UNIQUE (song_url),
			PRIMARY KEY (id)
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			uuid VARCHAR(64) NOT NULL,
			username VARCHAR(25) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password_hash VARCHAR(64) NOT NULL,
			role CHAR(1) NOT NULL,
			UNIQUE (username),
			UNIQUE (email),
			PRIMARY KEY (uuid)
		);`,
	}

	err := config.EnsurePQReady(prepareDBStatements)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.PQConn())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &pQDatastore{db}, nil
}
