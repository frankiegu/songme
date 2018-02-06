package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/controllers"
	"github.com/emre-demir/songme/datastore"
)

func main() {
	// Database
	db, err := datastore.NewDatastore(datastore.PostgreSQL, datastore.DefaultConfig)
	if err != nil {
		log.Fatal("[MAIN - DB]:", err)
	}
	defer db.Close()

	env := &env.Env{DB: db}

	// FileServer
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", controllers.Index)
	http.Handle("/signup", controllers.SignUp(env))

	// Serve
	port := "8080"
	if os.Getenv("ENV") == "PRODUCTION" {
		port = os.Getenv("PORT")
		if port == "" {
			log.Fatal("MAIN - PORT:", err)
		}
	}
	http.ListenAndServe(":"+port, nil)
}
