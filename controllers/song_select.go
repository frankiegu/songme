package controllers

import (
	"log"
	"net/http"

	"github.com/emre-demir/songme/common/env"
	"github.com/gorilla/mux"
)

// SelectSong transfers song from table 'candidate_song' to 'production_song'.
func SelectSong(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err := ev.DB.SelectSong(id)
		if err != nil {
			log.Println("[SelectSong]:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Song with id:", id, "transferred to table 'production_song'.")
		w.WriteHeader(http.StatusOK)
	})
}
