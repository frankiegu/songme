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
		http.NotFound(w, r)
	}
}
