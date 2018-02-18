package controllers

import (
	"net/http"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
)

type dashboardViewData struct {
	Songs  []*models.Song
	Errors map[string]string
}

// DashboardController renders dashboard view.
func DashboardController(ev *env.Vars) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			showDashboardView(w, r, ev)
		default:
			http.NotFound(w, r)
		}
	})
}

// showDashboardView fetch songs, and then renders admin/dashboard view.
func showDashboardView(w http.ResponseWriter, r *http.Request, ev *env.Vars) {
	dvd := dashboardViewData{
		Songs:  []*models.Song{},
		Errors: make(map[string]string),
	}

	songs, err := ev.DB.GetSongsFrom("candidate_song")
	if err != nil {
		dvd.Errors["database"] = "Could not fetch any song from database."
		RenderTemplate(w, "admin/dashboard", &dvd)
		return
	}

	dvd.Songs = songs
	RenderTemplate(w, "admin/dashboard", &dvd)
}
