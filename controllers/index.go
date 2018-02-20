package controllers

import (
	"net/http"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
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
