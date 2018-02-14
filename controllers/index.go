package controllers

import (
	"net/http"

	"github.com/emre-demir/songme/common/env"
)

// Index displays home page.
func Index(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := map[string]string{
			"PageTitle": "Songme",
		}
		switch r.Method {
		case "GET":
			RenderTemplate(w, "index/home", &p)
		case "POST":
			http.NotFound(w, r)
		}
	})
}
