package datastore

import (
	"database/sql"
	// Run postgresql driver.
	_ "github.com/lib/pq"

	"github.com/emre-demir/songme/models"
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
	INSERT INTO song (
		title, 
		author, 
		song_url, 
		image_url, 
		description, 
		recommended
	) 
	VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := pq.Exec(stmt, song.Title, song.Author, song.SongURL, song.ImageURL, song.Description, song.Recommended)
	return err
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
		`CREATE TABLE IF NOT EXISTS song (
			id SERIAL,
			title text NOT NULL,
			author text NOT NULL,
			song_url text NOT NULL,
			image_url text,
			description text,
			recommended BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			recommended_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			PRIMARY KEY (id)
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
