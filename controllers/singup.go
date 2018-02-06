package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
)

// indexPage contains information about index page.
type indexPage struct {
	Title        string
	ElementNames []string
	Fields       map[string]string
	Errors       map[string]string
}

// SignUp registers new subscriber
func SignUp(e *env.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusNotFound)
			RenderTemplate(w, "error/page", nil)
		case "POST":
			p := &indexPage{
				Title:        "Songme",
				ElementNames: []string{"name", "email"},
				Fields:       make(map[string]string),
				Errors:       make(map[string]string),
			}
			for _, element := range p.ElementNames {
				p.Fields[element] = r.FormValue(element)
			}
			signup(w, r, e, p)
		}
	})
}

// signup validates and handles posted singup form
func signup(w http.ResponseWriter, r *http.Request, e *env.Env, p *indexPage) {
	name := r.FormValue("name")
	email := r.FormValue("email")

	if strings.TrimSpace(name) == "" {
		p.Errors["name"] = "Please enter a valid name."
	}

	if strings.TrimSpace(email) == "" {
		p.Errors["email"] = "Please enter a valid email."
	}

	if len(p.Errors) > 0 {
		RenderTemplate(w, "index/home", &p)
		return
	}

	subscriber := &models.Subscriber{
		Name:  name,
		Email: email,
	}
	err := e.DB.Subscribe(subscriber)
	if err != nil {
		log.Println("[SignUp]:", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "error/page", nil)
		return
	}

	// TODO: later, make this method ready for production

	http.Redirect(w, r, "/", http.StatusSeeOther)

	log.Println("Subscribed:", name, email)
}
