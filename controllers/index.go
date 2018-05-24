package controllers

import (
	"net/http"

	"github.com/emredir/songme/common/env"
	"github.com/emredir/songme/models"
)

type homeViewData struct {
	Song *models.Song
}

// IndexController handles GET requests on path '/'.
// It shows home view with a randomly selected song.
func IndexController(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showHomeView(w, r, ev)
	})
}

// showHomeView fetches a random song, and renders home view
// with that song.
func showHomeView(w http.ResponseWriter, r *http.Request, ev *env.Vars) {
	s := ev.DB.GetRandomSong()
	hvd := homeViewData{Song: s}
	RenderTemplate(w, "index/home", &hvd)
}
