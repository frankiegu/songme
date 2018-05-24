package controllers

import (
	"net/http"

	"github.com/emredir/songme/common/env"
	"github.com/emredir/songme/models"
)

type dashboardViewData struct {
	Songs  []*models.Song
	Hidden map[string]bool
	Errors map[string]string
}

// DashboardController renders dashboard view.
func DashboardController(ev *env.Vars) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showDashboardView(w, r, ev)
	})
}

// showDashboardView fetches songs from table candidate_song, and then
// renders admin/dashboard view.
func showDashboardView(w http.ResponseWriter, r *http.Request, ev *env.Vars) {
	dvd := dashboardViewData{
		Songs:  []*models.Song{},
		Hidden: make(map[string]bool),
		Errors: make(map[string]string),
	}

	songs, err := ev.DB.GetSongsFrom("candidate_song")
	if err != nil {
		dvd.Errors["database"] = "database-error: Could not fetch any song from database."
		RenderTemplate(w, "admin/dashboard", &dvd)
		return
	}

	dvd.Songs = songs
	RenderTemplate(w, "admin/dashboard", &dvd)
}
