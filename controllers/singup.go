package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
)

// SignUp registers new subscriber
func SignUp(e *env.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		case "POST":
			signup(w, r, e)
		}
	})
}

// signup validates and handles posted singup form
func signup(w http.ResponseWriter, r *http.Request, e *env.Env) {
	name := r.FormValue("name")
	email := r.FormValue("email")

	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	subscriber := &models.Subscriber{
		Name:  name,
		Email: email,
	}
	err := e.DB.Subscribe(subscriber)
	if err != nil {
		log.Println("[SignUp]:", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// TODO: later, make this method ready for production

	http.Redirect(w, r, "/", http.StatusSeeOther)

	log.Println("Subscribed:", name, email)
}
