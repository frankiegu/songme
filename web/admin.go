package web

import (
	"net/http"

	"github.com/emredir/songme/models"
)

type dashboardViewData struct {
	Songs  []*models.Song
	Hidden map[string]bool
	Errors map[string]string
}

// AdminHandler defines admin specific controllers.
type AdminHandler struct {
	songStore models.SongStore
}

// Dashboard renders the dashboard page.
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	dvd := dashboardViewData{
		Songs:  []*models.Song{},
		Hidden: make(map[string]bool),
		Errors: make(map[string]string),
	}

	songs, err := h.songStore.Candidates()
	if err != nil {
		dvd.Errors["database"] = "database-error: Could not fetch any song from database."
		RenderTemplate(w, "admin/dashboard", &dvd)
		return
	}

	dvd.Songs = songs
	RenderTemplate(w, "admin/dashboard", &dvd)
}
