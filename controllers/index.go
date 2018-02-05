package controllers

import (
	"net/http"
)

// Index displays home page.
func Index(w http.ResponseWriter, r *http.Request) {
	page := map[string]string{
		"Title": "Songme",
	}

	switch r.Method {
	case "GET":
		RenderTemplate(w, "index/home", &page)
	case "POST":
		// TODO: make sure http.NotFound is appropriate.
		// TODO: display a custom 'not found' page.
		http.NotFound(w, r)
	}
}
