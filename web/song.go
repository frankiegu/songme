package web

import (
	"net/http"
	"strconv"

	"github.com/emredir/songme/models"

	"github.com/gorilla/mux"
)

// SongInteractor is used to interact with the database for song related tasks.
type SongInteractor interface {
	models.SongStore
}

// SongHandler defines song specific controllers.
type SongHandler struct {
	songInteractor    SongInteractor
	descriptionLength int
	songsPerPage      int
}

// New renders the form for creating a new song.
func (h *SongHandler) New(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	view.Render(w, "song/recommend")
}

// Create creates a new song.
func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	title := view.FormValue("title", true)
	artist := view.FormValue("artist", true)
	songURL := view.FormValue("songURL", true)
	imageURL := view.FormValue("imageURL", false)
	description := view.FormValue("description", false)

	if view.HasError() {
		view.Render(w, "song/recommend")
		return
	}

	if len(description) > h.descriptionLength {
		view.InsertFlashError("Descriptions must be at most ", h.descriptionLength, " characters long")
		view.Render(w, "song/recommend")
		return
	}

	s := &models.Song{
		Title:       title,
		Artist:      artist,
		SongURL:     songURL,
		ImageURL:    &imageURL,
		Description: &description,
	}
	err := h.songInteractor.Create(s)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "song/recommend")
		return
	}

	view.InsertSong(s)
	view.Render(w, "song/success")
}

// Songs returns all songs that are currently confirmed.
func (h *SongHandler) Songs(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	page := 1
	vars := mux.Vars(r)
	val, ok := vars["page"]
	if ok {
		var err error
		page, err = strconv.Atoi(val)
		if err != nil {
			view.InsertFlashError("Page ", val, " cannot be found")
			view.Render(w, "song/all")
			return
		}
	}

	offset := (page - 1) * h.songsPerPage
	songs, totalCount, err := h.songInteractor.All(true, h.songsPerPage, offset)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "song/all")
		return
	}

	pagination := newPagination("songs", totalCount, h.songsPerPage, 10, page)
	view.InsertSongs(songs)
	view.InsertPagination(pagination)
	view.Render(w, "song/all")
}
