package main

import (
	"log"
	"net/http"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/controllers"
	"github.com/emre-demir/songme/datastore"
)

func main() {
	db, err := datastore.NewDatastore(datastore.PostgreSQL, datastore.DefaultConfig)
	if err != nil {
		log.Fatal("[MAIN - DB]:", err)
	}
	defer db.Close()

	env := &env.Env{DB: db}

	// FileServer
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", controllers.Index)
	http.Handle("/signup", controllers.SignUp(env))

	http.ListenAndServe(":8080", nil)
}
