package web

import (
	"net/http"

	"github.com/emredir/songme/models"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type songAddViewData struct {
	Form *models.Form
}

type songSuccessViewData struct {
	Song *models.Song
}

type homeViewData struct {
	Song *models.Song
}

// SongHandler defines song specific controllers.
type SongHandler struct {
	songStore models.SongStore
}

// Index renders the home page with a randomly selected song.
func (h *SongHandler) Index(w http.ResponseWriter, r *http.Request) {
	s, _ := h.songStore.GetRandom()
	RenderTemplate(w, "index/home", &homeViewData{Song: s})
}

// New renders the form for creating a new song.
func (h *SongHandler) New(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "song/add", nil)
}

// Create creates a new song.
func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	f := models.NewForm([]string{"title", "author", "songURL", "imageURL", "description"})

	f.Populate(r)
	f.CheckRequiredFields([]string{"title", "author", "songURL"})
	ok := f.IsValid()
	if !ok {
		RenderTemplate(w, "song/add", &songAddViewData{Form: f})
		return
	}

	s := &models.Song{
		Title:       r.FormValue("title"),
		Author:      r.FormValue("author"),
		SongURL:     r.FormValue("songURL"),
		ImageURL:    r.FormValue("imageURL"),
		Description: r.FormValue("description"),
	}
	err := h.songStore.Create(s)
	if err == nil {
		RenderTemplate(w, "song/success", &songSuccessViewData{Song: s})
		return
	}

	// Check for the type of postgresql error.
	if pqError, ok := err.(*pq.Error); ok {
		// If error is because of unique violation, then inform user.
		if pqError.Code == "23505" {
			f.Errors["alreadyExists"] = "Opps! Sorry, song already exists."
			RenderTemplate(w, "song/add", &songAddViewData{Form: f})
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "error/page", nil)
}

// Confirm transfers song from table 'candidate_song' to 'production_song'.
func (h *SongHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.songStore.Confirm(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Candidates renders all candidate songs.
func (h *SongHandler) Candidates(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songStore.Candidates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["Songs"] = songs

	RenderTemplate(w, "song/rows", data)
}

// Productions renders all songs in production.
func (h *SongHandler) Productions(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songStore.Productions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["Songs"] = songs
	data["Hidden"] = map[string]bool{"selectButton": true}

	RenderTemplate(w, "song/rows", data)
}

// DeleteCandidate deletes song from table 'candidate_song'.
func (h *SongHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.songStore.DeleteCandidate(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProduction deletes song from table 'production_song'.
func (h *SongHandler) DeleteProduction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.songStore.DeleteProduction(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
