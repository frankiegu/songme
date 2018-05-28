package web

import (
	"log"
	"net/http"

	"github.com/emredir/songme/common/auth"
	"github.com/emredir/songme/common/utility"
	"github.com/emredir/songme/models"
)

type loginViewData struct {
	Form *models.Form
}

// AuthHandler defines authentication specific controllers.
type AuthHandler struct {
	userStore models.UserStore
}

// RenderLogin renders the login page.
func (h *AuthHandler) RenderLogin(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "session/login", nil)
}

// Login handles login process.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	f := models.NewForm([]string{"username", "password"})

	f.Populate(r)
	f.CheckRequiredFields([]string{"username", "password"})
	ok := f.IsValid()
	if !ok {
		RenderTemplate(w, "session/login", &loginViewData{Form: f})
		return
	}

	username := f.Inputs["username"]
	password := f.Inputs["password"]

	u, err := h.userStore.ByUsername(username)
	if err != nil {
		f.Errors["login"] = "Please check your username and password."
		RenderTemplate(w, "session/login", &loginViewData{Form: f})
		return
	}

	if utility.SHA256String(password) != u.PasswordHash {
		f.Errors["login"] = "Please check your username and password."
		RenderTemplate(w, "session/login", &loginViewData{Form: f})
		return
	}

	log.Println("User logged in:", u.Username)

	err = auth.SetSessionCookie(w, u)
	if err != nil {
		f.Errors["login"] = "Opps! Something went wrong. Please try again."
		RenderTemplate(w, "session/login", &loginViewData{Form: f})
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

// Logout logouts user.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
