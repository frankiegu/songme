package web

import (
	"html/template"
	"log"
	"net/http"
)

var (
	templates = template.Must(template.New("start").Funcs(template.FuncMap{
		"toHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).ParseGlob("templates/**/*.html"))
)

// RenderTemplate renders given template
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("[RenderTemplate]:", err)
		http.Error(w, "Opps! Something went wrong. We are going to take care of it.", http.StatusInternalServerError)
		return
	}
}
