package controllers

import (
	"html/template"
	"log"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("templates/**/*.html"))
)

// RenderTemplate renders given template
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("[RenderTemplate]:", err)
		http.Error(w, "Opps! There is something wrong. We are going to take care of it.", http.StatusInternalServerError)
		return
	}
}
