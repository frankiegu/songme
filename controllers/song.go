package controllers

import (
	"log"
	"net/http"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
	"github.com/lib/pq"
)

type songAddViewData struct {
	Form *models.Form
}

type songSuccessViewData struct {
	Song *models.Song
}

// AddSongController handles GET&POST requests on path /add-song.
// Renders song/add view on a successful GET request.
// If the case is POST then it handles persisting song to database.
func AddSongController(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			RenderTemplate(w, "song/add", nil)
		case "POST":
			handleAddSong(w, r, ev)
		}
	})
}

// handleAddSong processes the given request.
// First, it populates form values, then it makes sure
// required fields are filled. Finally, inserts new song into database.
func handleAddSong(w http.ResponseWriter, r *http.Request, ev *env.Vars) {
	f := models.NewForm([]string{"title", "author", "songURL", "imageURL", "description"})

	f.Populate(r)
	f.CheckRequiredFields([]string{"title", "author", "songURL"})

	ok := f.IsValid()
	if !ok {
		vd := songAddViewData{Form: f}
		RenderTemplate(w, "song/add", &vd)
		return
	}

	s := models.NewSong(
		r.FormValue("title"),
		r.FormValue("author"),
		r.FormValue("songURL"),
		r.FormValue("imageURL"),
		r.FormValue("description"),
	)

	err := ev.DB.CreateSong(s)
	if err == nil {
		log.Println("Song successfully added:", r.FormValue("title"), r.FormValue("author"))
		vd := songSuccessViewData{Song: s}
		RenderTemplate(w, "song/success", &vd)
		return
	}

	// Check for the type of postgresql error.
	if pqError, ok := err.(*pq.Error); ok {
		// If error is because of unique violation, then inform user.
		if pqError.Code == "23505" {
			f.Errors["alreadyExists"] = "Opps! Sorry, song already exists."
			vd := songAddViewData{Form: f}
			RenderTemplate(w, "song/add", &vd)
			return
		}
	}

	log.Println("[handleAddSong:]", err)
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "error/page", nil)
}
