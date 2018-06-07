package web

import (
	"net/http"
	"strconv"

	"github.com/emredir/songme/models"
	"github.com/gorilla/mux"
)

// AdminHandler defines admin specific controllers.
type AdminHandler struct {
	userStore    models.UserStore
	songStore    models.SongStore
	songsPerPage int
	usersPerPage int
}

// Dashboard renders the dashboard template with candidate songs.
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	page := 1
	vars := mux.Vars(r)
	val, ok := vars["page"]
	if ok {
		var err error
		page, err = strconv.Atoi(val)
		if err != nil {
			view.InsertFlashError("Page ", val, " cannot be found")
			view.Render(w, "admin/dashboard")
			return
		}
	}

	offset := (page - 1) * h.songsPerPage
	songs, totalCount, err := h.songStore.All(false, h.songsPerPage, offset)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "admin/dashboard")
		return
	}

	pagination := newPagination("admin/dashboard", totalCount, h.songsPerPage, 10, page)
	view.InsertSongs(songs)
	view.InsertPagination(pagination)
	view.InsertOther("songs", totalCount)
	view.Render(w, "admin/dashboard")
}

// Productions renders dashboard template with production songs.
func (h *AdminHandler) Productions(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	page := 1
	vars := mux.Vars(r)
	val, ok := vars["page"]
	if ok {
		var err error
		page, err = strconv.Atoi(val)
		if err != nil {
			view.InsertFlashError("Page ", val, " cannot be found")
			view.Render(w, "admin/dashboard")
			return
		}
	}

	offset := (page - 1) * h.songsPerPage
	songs, totalCount, err := h.songStore.All(true, h.songsPerPage, offset)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "admin/dashboard")
		return
	}

	pagination := newPagination("admin/dashboard/productions", totalCount, h.songsPerPage, 10, page)
	view.InsertSongs(songs)
	view.InsertPagination(pagination)
	view.InsertOther("songs", totalCount)
	view.Render(w, "admin/dashboard")
}

// Users renders dashboard template with users.
func (h *AdminHandler) Users(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	page := 1
	vars := mux.Vars(r)
	val, ok := vars["page"]
	if ok {
		var err error
		page, err = strconv.Atoi(val)
		if err != nil {
			view.InsertFlashError("Page ", val, " cannot be found")
			view.Render(w, "admin/dashboard")
			return
		}
	}

	offset := (page - 1) * h.usersPerPage
	users, totalCount, err := h.userStore.All(h.usersPerPage, offset)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "admin/dashboard")
		return
	}

	pagination := newPagination("admin/dashboard/users", totalCount, h.usersPerPage, 10, page)
	view.InsertUsers(users)
	view.InsertPagination(pagination)
	view.InsertOther("users", totalCount)
	view.Render(w, "admin/dashboard")
}

// ConfirmSong updates song's confirmation status to true.
func (h *AdminHandler) ConfirmSong(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		view.InsertFlashError("ID ", id, " cannot be found")
		view.Render(w, "admin/dashboard")
		return
	}

	err := h.songStore.Confirm(id)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "admin/dashboard")
		return
	}

	view.InsertFlash("Song with ID ", id, " successfully confirmed")
	view.Render(w, "admin/dashboard")
}

// DeleteSong deletes a song from the database.
func (h *AdminHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		view.InsertFlashError("ID ", id, " cannot be found")
		view.Render(w, "admin/dashboard")
		return
	}

	err := h.songStore.Delete(id)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "admin/dashboard")
		return
	}

	view.InsertFlash("Song with ID ", id, " successfully deleted")
	view.Render(w, "admin/dashboard")
}
