package web

import (
	"net/http"

	"github.com/emredir/songme/internal/context"
	"github.com/emredir/songme/internal/cookie"
	"github.com/emredir/songme/models"
)

// AuthInteractor defines the interface used to interact with the database for authentication.
type AuthInteractor interface {
	Signup(email, username, password string) (*models.User, error)
	Signin(username, password string) (*models.User, error)
	UpdateEmail(old, new, password string) error
	UpdatePassword(email, oldPassword, newPassword string) error
}

// AuthHandler defines authentication specific controllers.
type AuthHandler struct {
	AuthInteractor AuthInteractor
	UsernameLength int
	PasswordLength int
}

/*

	Controllers

*/

// RenderSignup renders the sign up page.
func (h *AuthHandler) RenderSignup(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	view.Render(w, "auth/signup")
}

// Signup handles sign up process.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	email := view.FormValue("email", true)
	username := view.FormValue("username", true)
	password := view.FormValue("password", true)
	password2 := view.FormValue("password2", true)

	if view.HasError() {
		view.Render(w, "auth/signup")
		return
	}

	if len(username) < h.UsernameLength {
		view.InsertFlashError("Usernames must be at least ", h.UsernameLength, " characters long")
		view.Render(w, "auth/signup")
		return
	}
	if len(password) < h.PasswordLength || len(password2) < h.PasswordLength {
		view.InsertFlashError("Passwords must be at least ", h.PasswordLength, " characters long")
		view.Render(w, "auth/signup")
		return
	}
	if password != password2 {
		view.InsertFlashError("Passwords must be matched")
		view.Render(w, "auth/signup")
		return
	}

	_, err := h.AuthInteractor.Signup(email, username, password)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "auth/signup")
		return
	}

	http.Redirect(w, r, "/signin", http.StatusFound)
}

// RenderSignin renders the sign in page.
func (h *AuthHandler) RenderSignin(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	view.Render(w, "auth/signin")
}

// Signin handles sign in process.
func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	username := view.FormValue("username", true)
	password := view.FormValue("password", true)

	if view.HasError() {
		view.Render(w, "auth/signin")
		return
	}

	user, err := h.AuthInteractor.Signin(username, password)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "auth/signin")
		return
	}

	err = cookie.Set(w, user.ID)
	if err != nil {
		view.InsertFlashError("Opps! Something went wrong. Error: ", err.Error())
		view.Render(w, "auth/signin")
		return
	}

	http.Redirect(w, r, "/user/"+user.Username, http.StatusFound)
}

// Logout logouts user.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie.Clear(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// RenderUpdatePassword renders the password update template.
func (h *AuthHandler) RenderUpdatePassword(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	view.Render(w, "auth/password")
}

// UpdatePassword handles password changing process.
func (h *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	email := user.Email
	oldPassword := view.FormValue("oldPassword", true)
	newPassword := view.FormValue("newPassword", true)
	confirmPassword := view.FormValue("confirmPassword", true)

	if view.HasError() {
		view.Render(w, "auth/password")
		return
	}

	if len(newPassword) < h.PasswordLength || len(confirmPassword) < h.PasswordLength {
		view.InsertFlashError("Passwords must be at least ", h.PasswordLength, " characters long")
		view.Render(w, "auth/password")
		return
	}
	if newPassword != confirmPassword {
		view.InsertFlashError("Passwords must be matched")
		view.Render(w, "auth/password")
		return
	}

	err := h.AuthInteractor.UpdatePassword(email, oldPassword, newPassword)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "auth/password")
		return
	}

	view.InsertFlash("Your password succesfully changed")
	view.Render(w, "auth/password")
}

// RenderUpdateEmail renders the email update template.
func (h *AuthHandler) RenderUpdateEmail(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)
	view.Render(w, "auth/email")
}

// UpdateEmail handles email changing process.
func (h *AuthHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	view := NewView(r)

	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	oldEmail := user.Email
	newEmail := view.FormValue("newEmail", true)
	password := view.FormValue("password", true)

	if view.HasError() {
		view.Render(w, "auth/email")
		return
	}

	err := h.AuthInteractor.UpdateEmail(oldEmail, newEmail, password)
	if err != nil {
		view.InsertFlashError(err.Error())
		view.Render(w, "auth/email")
		return
	}

	view.InsertFlash("Your email succesfully changed to ", newEmail)
	view.Render(w, "auth/email")
}
