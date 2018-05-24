package controllers

import (
	"log"
	"net/http"

	"github.com/emredir/songme/common/env"
	"github.com/gorilla/mux"
)

// DeleteProductionSong deletes song from table 'production_song'.
func DeleteProductionSong(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err := ev.DB.DeleteProductionSong(id)
		if err != nil {
			log.Println("[DeleteProductionSong]:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Song with id:", id, "is deleted from table 'production_song'.")
		w.WriteHeader(http.StatusOK)
	})
}

// DeleteCandidateSong deletes song from table 'candidate_song'.
func DeleteCandidateSong(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err := ev.DB.DeleteCandidateSong(id)
		if err != nil {
			log.Println("[DeleteCandidateSong]:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Song with id:", id, "is deleted from table 'candidate_song'.")
		w.WriteHeader(http.StatusOK)
	})
}
