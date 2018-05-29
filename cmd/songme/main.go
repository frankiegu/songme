package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/emredir/songme"
	"github.com/emredir/songme/databases/psql"
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

	userStore := psql.UserStore{DB: db}
	songStore := psql.SongStore{DB: db}
	server := web.NewServer(&userStore, &songStore)

	// Serve
	log.Fatal(http.ListenAndServe(songme.GetConfig().Port, server))
}
