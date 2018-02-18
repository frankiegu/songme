package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/controllers"
	"github.com/emre-demir/songme/datastore"
	"github.com/emre-demir/songme/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Database
	db, err := datastore.NewDatastore(datastore.PostgreSQL, datastore.DefaultConfig)
	if err != nil {
		log.Fatal("[MAIN - DB]:", err)
	}
	defer db.Close()

	ev := &env.Vars{DB: db}

	r := mux.NewRouter()

	// Routes
	r.Handle("/", controllers.IndexController(ev)).Methods("GET")
	r.Handle("/add-song", controllers.AddSongController(ev)).Methods("GET", "POST")
	r.Handle("/login", controllers.LoginController(ev)).Methods("GET", "POST")
	r.Handle("/admin/dashboard", middleware.Authorize(controllers.DashboardController(ev), ev)).Methods("GET")

	// FileServer
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/", middleware.PanicRecovery(r))

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
