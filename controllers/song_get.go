package controllers

import (
	"log"
	"net/http"

	"github.com/emredir/songme/common/env"
)

// GetCandidateSongs fetches all songs from table 'candidate_song'.
func GetCandidateSongs(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchSongsFrom(w, ev, "candidate_song")
	})
}

// GetProductionSongs fetches all songs from table 'production_song'.
func GetProductionSongs(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchSongsFrom(w, ev, "production_song")
	})
}

// fetchSongsFrom fetches all songs from given table.
// After fetching all songs, it writes them into template song/rows
// and then sends that as a response.
func fetchSongsFrom(w http.ResponseWriter, ev *env.Vars, table string) {
	songs, err := ev.DB.GetSongsFrom(table)
	if err != nil {
		log.Println("[fetchSongsFrom]:", table, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["Songs"] = songs
	if table == "production_song" {
		data["Hidden"] = map[string]bool{"selectButton": true}
	}

	RenderTemplate(w, "song/rows", data)
}
