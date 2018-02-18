package controllers

import (
	"log"
	"net/http"

	"github.com/emre-demir/songme/common/auth"
	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
)

type loginViewData struct {
	Form *models.Form
}

// LoginController handles requests for signing in.
func LoginController(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			RenderTemplate(w, "session/login", nil)
		case "POST":
			handleLogin(w, r, ev)
		}
	})
}

// handleLogin handles login process.
func handleLogin(w http.ResponseWriter, r *http.Request, ev *env.Vars) {
	f := models.NewForm([]string{"username", "password"})

	f.Populate(r)
	f.CheckRequiredFields([]string{"username", "password"})
	ok := f.IsValid()
	if !ok {
		vd := loginViewData{Form: f}
		RenderTemplate(w, "session/login", &vd)
		return
	}

	u, err := auth.LoginUserWith(f.Inputs["username"], f.Inputs["password"], ev)
	if err != nil {
		f.Errors["login"] = "Please check your username and password."
		vd := loginViewData{Form: f}
		RenderTemplate(w, "session/login", &vd)
		return
	}

	err = auth.SetSessionCookie(w, u)
	if err != nil {
		f.Errors["login"] = "Opps! Something went wrong. Please try again."
		vd := loginViewData{Form: f}
		RenderTemplate(w, "session/login", &vd)
		return
	}

	log.Println("login successful for user:", r.FormValue("username"))
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}
