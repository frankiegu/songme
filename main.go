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

	// Routes
	// Main router
	router := mux.NewRouter()
	router.Handle("/", controllers.IndexController(ev)).Methods("GET")
	router.Handle("/add", controllers.AddSongController(ev)).Methods("GET", "POST")
	router.Handle("/login", controllers.LoginController(ev)).Methods("GET", "POST")
	router.Handle("/logout", controllers.LogoutController()).Methods("POST")

	// Admin router
	adminRouter := mux.NewRouter()
	adminRouter.Handle("/admin/dashboard", controllers.DashboardController(ev)).Methods("GET")

	// Songs router
	songsRouter := mux.NewRouter()
	songsRouter.Handle("/songs/{id}", controllers.SelectSong(ev)).Methods("PUT")
	songsRouter.Handle("/songs/candidate", controllers.GetCandidateSongs(ev)).Methods("GET")
	songsRouter.Handle("/songs/production", controllers.GetProductionSongs(ev)).Methods("GET")
	songsRouter.Handle("/songs/candidate/{id}", controllers.DeleteCandidateSong(ev)).Methods("DELETE")
	songsRouter.Handle("/songs/production/{id}", controllers.DeleteProductionSong(ev)).Methods("DELETE")

	// FileServer
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Authorize admin router
	router.PathPrefix("/admin").Handler(
		middleware.Authorize(adminRouter, ev),
	)

	// Authorize songs router
	router.PathPrefix("/songs").Handler(
		middleware.Authorize(songsRouter, ev),
	)

	// Recover panics on main router
	http.Handle("/", middleware.PanicRecovery(router))

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
