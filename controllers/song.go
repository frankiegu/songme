package controllers

import (
	"net/http"
	"strings"

	"github.com/emre-demir/songme/common/env"
)

// form contains basic information about forms on html page.
type form struct {
	PageTitle      string
	RequiredInputs []string
	Inputs         map[string]string
	Errors         map[string]string
}

// isFilled checks required input fields in order to
// validate whether they are filled or not.
func (f *form) isFilled(r *http.Request) bool {
	for _, input := range f.RequiredInputs {
		val := r.FormValue(input)
		f.Inputs[input] = val
		if strings.TrimSpace(val) == "" {
			f.Errors[input] = "Please enter a valid " + input + "."
		}
	}
	return len(f.Errors) == 0
}

// AddSong renders song/add.html page.
func AddSong(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := &form{
			PageTitle:      "Songme",
			RequiredInputs: []string{"title", "author", "song-url"},
			Inputs:         make(map[string]string),
			Errors:         make(map[string]string),
		}
		switch r.Method {
		case "GET":
			RenderTemplate(w, "song/add", f)
		case "POST":
			addSong(w, r, ev, f)
		}
	})
}

// addSong reads posted form and adds new song
func addSong(w http.ResponseWriter, r *http.Request, ev *env.Vars, f *form) {
	http.NotFound(w, r)
}
