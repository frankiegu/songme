package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/emredir/songme"
	"github.com/emredir/songme/databases/psql"
	"github.com/emredir/songme/interactor"
	"github.com/emredir/songme/web"
)

func main() {
	// Connect to database
	db, err := sql.Open("postgres", songme.GetConfig().DatabaseURL)
	if err != nil {
		log.Fatal("[MAIN - DB]:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("[MAIN - DB]:", err)
	}

	store := web.Store{
		User: &psql.UserStore{DB: db},
		Role: &psql.RoleStore{DB: db},
		Song: &psql.SongStore{DB: db},
	}
	interactor := web.Interactor{
		Auth: &interactor.Auth{
			UserStore: store.User,
			RoleStore: store.Role,
		},
		Song: &interactor.Song{
			SongStore: store.Song,
		},
	}
	server := web.NewServer(store, interactor)

	// Serve
	log.Fatal(http.ListenAndServe(songme.GetConfig().Port, server))
}
