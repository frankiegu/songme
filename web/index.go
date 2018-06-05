package web

import (
	"net/http"

	"github.com/emredir/songme/models"
)

// MainHandler defines the main controllers such as home, contact us, about us etc.
type MainHandler struct {
	songStore models.SongStore
}

// Home renders the home page with a randomly selected song.
func (h *MainHandler) Home(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	s, err := h.songStore.GetRandom()
	if err != nil {
		view.InsertFlashError(err.Error())
	} else {
		view.InsertSong(s)
	}
	view.Render(w, "home/home")
}
