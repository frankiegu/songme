package web

import (
	"net/http"
	"strconv"

	"github.com/emredir/songme/models"
	"github.com/gorilla/mux"
)

// UserHandler defines user specific controllers.
type UserHandler struct {
	userStore    models.UserStore
	songStore    models.SongStore
	songsPerPage int
}

// Profile renders profile page for the current user.
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	vars := mux.Vars(r)
	username := vars["username"]

	user, err := h.userStore.ByUsername(username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	page := 1
	val, ok := vars["page"]
	if ok {
		var err error
		page, err = strconv.Atoi(val)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	}

	offset := (page - 1) * h.songsPerPage
	songs, totalCount, err := h.songStore.UserSongs(user.ID, h.songsPerPage, offset)
	if err != nil {
		view.InsertFlashError("Sorry, it seems we can't find what you're looking for")
		view.Render(w, "user/profile")
		return
	}

	pagination := newPagination("user/"+username, totalCount, h.songsPerPage, 10, page)
	view.InsertUser(user)
	view.InsertSongs(songs)
	view.InsertPagination(pagination)
	view.InsertOther("recommendations", totalCount)
	view.Render(w, "user/profile")
}
