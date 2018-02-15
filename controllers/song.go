package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
	"github.com/lib/pq"
)

// form contains basic information about forms on html page.
type form struct {
	InputNames []string
	Inputs     map[string]string
	Errors     map[string]string
}

// populate extracts form values from http.Request
// and fills it's Inputs.
func (f *form) populate(r *http.Request) {
	for _, input := range f.InputNames {
		f.Inputs[input] = r.FormValue(input)
	}
}

// checkRequiredFields ensures that required fields
// filled up properly
func (f *form) checkRequiredFields(fieldNames []string) {
	for _, field := range fieldNames {
		if strings.TrimSpace(f.Inputs[field]) == "" {
			f.Errors[field] = "Please enter a valid " + field + "."
		}
	}
}

// isValid checks whether there is an error or not.
func (f *form) isValid() bool {
	return len(f.Errors) == 0
}

// AddSongController handles displaying and processing song forms.
func AddSongController(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := &form{
			InputNames: []string{"title", "author", "songURL", "imageURL", "description"},
			Inputs:     make(map[string]string),
			Errors:     make(map[string]string),
		}
		switch r.Method {
		case "GET":
			RenderTemplate(w, "song/add", f)
		case "POST":
			processAddSong(w, r, ev, f)
		}
	})
}

// processAddSong processes the given request.
// First, it populates form values, then it makes sure
// required fields are filled. Finally, inserts newly added song
// into the database.
func processAddSong(w http.ResponseWriter, r *http.Request, ev *env.Vars, f *form) {
	f.populate(r)
	f.checkRequiredFields([]string{"title", "author", "songURL"})

	ok := f.isValid()
	if !ok {
		RenderTemplate(w, "song/add", f)
		return
	}

	song := models.NewSong(
		r.FormValue("title"),
		r.FormValue("author"),
		r.FormValue("songURL"),
		r.FormValue("imageURL"),
		r.FormValue("description"),
	)

	err := ev.DB.CreateSong(song)
	if err == nil {
		log.Println("Song successfully added:", r.FormValue("title"), r.FormValue("author"))
		RenderTemplate(w, "song/success", f)
		return
	}

	// Check for the type of postgresql error.
	if pqError, ok := err.(*pq.Error); ok {
		// If error is because of unique violation, then inform user.
		if pqError.Code == "23505" {
			f.Errors["alreadyExists"] = "Opps. Sorry, song already exists."
			RenderTemplate(w, "song/add", f)
			return
		}
	}

	log.Println("[addSong:]", err)
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "error/page", nil)
}
