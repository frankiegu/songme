package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/emre-demir/songme/common/env"
	"github.com/emre-demir/songme/models"
)

// sform contains basic information about subscription form.
type sform struct {
	PageTitle      string
	ComponentNames []string
	Components     map[string]string
	Errors         map[string]string
}

// Index displays home page.
func Index(ev *env.Vars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := &sform{
			PageTitle:      "Songme",
			ComponentNames: []string{"name", "email"},
			Components:     make(map[string]string),
			Errors:         make(map[string]string),
		}
		switch r.Method {
		case "GET":
			RenderTemplate(w, "index/home", s)
		case "POST":
			process(w, r, ev, s)
		}
	})
}

// process processes the given request. First, it
// extracts form values, then it makes sure
// components needed are filled.
func process(w http.ResponseWriter, r *http.Request, ev *env.Vars, s *sform) {
	for _, cname := range s.ComponentNames {
		val := r.FormValue(cname)
		s.Components[cname] = val
		if strings.TrimSpace(val) == "" {
			s.Errors[cname] = "Please enter a valid " + cname
		}
	}
	if len(s.Errors) > 0 {
		RenderTemplate(w, "index/home", s)
		w.Write([]byte("<script>document.getElementById('subscription-btn').click();</script>"))
		return
	}
	subscribe(w, r, ev, s)
}

// subscribe registers new subscriber.
func subscribe(w http.ResponseWriter, r *http.Request, ev *env.Vars, s *sform) {
	subs := models.NewSubscriber(r.FormValue("name"), r.FormValue("email"))
	err := ev.DB.Subscribe(subs)
	if err != nil {
		// TODO: Check if email already exists, if so, shows 'index/home'.
		log.Println("[subscribe]:", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "error/page", nil)
		return
	}

	// TODO: Show a 'subscribed successful' page.
	log.Println("Subscribed:", r.FormValue("name"), r.FormValue("email"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
