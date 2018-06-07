package web

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/emredir/songme/internal/context"
	"github.com/emredir/songme/models"
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.New("t").Funcs(template.FuncMap{
		"toHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
		"repeating": func(n int) []int {
			return make([]int, n)
		},
		"sum": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}).ParseFiles("templates/base.html"))

	templates = template.Must(templates.ParseGlob("templates/**/*.html"))
}

// NewView returns a new view object.
func NewView(r *http.Request) *View {
	view := &View{
		r:    r,
		Data: make(map[string]interface{}),
	}
	view.InserCurrentUser(context.User(r.Context()))
	return view
}

// View will hold the data that will be inserted in our views.
type View struct {
	r    *http.Request
	Data map[string]interface{}
}

// Render renders the view.
func (v *View) Render(w http.ResponseWriter, name string) {
	err := templates.ExecuteTemplate(w, name, v.Data)
	if err != nil {
		// TODO: Use a logger interface
		log.Println("[RenderTemplate]:", err)
		http.Error(w, "Opps! Something went wrong. We are going to take care of it.", http.StatusInternalServerError)
		return
	}
}

// HasError tells whether the view contains any error.
func (v *View) HasError() bool {
	_, e1 := v.Data["Error"].(map[string]string)
	_, e2 := v.Data["FlashError"].([]string)
	return e1 || e2
}

// FormValue wraps http.Request.FormValue.
// It extracts the form value and generates an error if required is true.
func (v *View) FormValue(input string, required bool) string {
	value := strings.TrimSpace(v.r.FormValue(input))
	if required && value == "" {
		v.InsertError(input, "Please enter a valid input")
	} else {
		v.InsertForm(input, value)
	}
	return value
}

// InserCurrentUser inserts current user into view.
func (v *View) InserCurrentUser(user *models.User) {
	v.Data["CurrentUser"] = user
}

// InsertUser inserts user into view.
func (v *View) InsertUser(user *models.User) {
	v.Data["User"] = user
}

// InsertUsers inserts slice of users into view.
func (v *View) InsertUsers(users []*models.User) {
	v.Data["Users"] = users
}

// InsertFlash inserts a flash message into view.
func (v *View) InsertFlash(a ...interface{}) {
	message := fmt.Sprint(a...)
	flash, ok := v.Data["Flash"].([]string)
	if !ok {
		v.Data["Flash"] = []string{message}
		return
	}
	v.Data["Flash"] = append(flash, message)
}

// InsertFlashError flashes an error message into view.
func (v *View) InsertFlashError(a ...interface{}) {
	message := fmt.Sprint(a...)
	flash, ok := v.Data["FlashError"].([]string)
	if !ok {
		v.Data["FlashError"] = []string{message}
		return
	}
	v.Data["FlashError"] = append(flash, message)
}

// InsertError inserts an error message into view.
func (v *View) InsertError(name, value string) {
	errors, ok := v.Data["Error"].(map[string]string)
	if !ok {
		v.Data["Error"] = map[string]string{name: value}
		return
	}
	errors[name] = value
}

// InsertForm inserts an input into view's form.
func (v *View) InsertForm(input, value string) {
	form, ok := v.Data["Form"].(map[string]string)
	if !ok {
		v.Data["Form"] = map[string]string{input: value}
		return
	}
	form[input] = value
}

// InsertSong inserts a song into view.
func (v *View) InsertSong(song *models.Song) {
	v.Data["Song"] = song
}

// InsertSongs inserts slice of songs into view.
func (v *View) InsertSongs(songs []*models.Song) {
	v.Data["Songs"] = songs
}

// InsertPagination inserts pagination info into view.
func (v *View) InsertPagination(pagination *Pagination) {
	v.Data["Pagination"] = pagination
}

// InsertOther inserts another key, value pair into view.
func (v *View) InsertOther(key string, value interface{}) {
	other, ok := v.Data["Other"].(map[string]interface{})
	if !ok {
		v.Data["Other"] = map[string]interface{}{key: value}
		return
	}
	other[key] = value
}

/*

	Pagination

*/

// Pagination will give pagination feature to our views.
type Pagination struct {
	Base       string
	Page       int
	Start      int
	BarSize    int
	TotalPages int
}

func newPagination(base string, total, perSize, barSize, page int) *Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(perSize)))
	start := 1

	if totalPages < barSize {
		barSize = totalPages
	} else {
		if page > 2 {
			start = page - 2
		}
		if start+barSize-1 >= totalPages {
			start = totalPages - barSize + 1
		}
	}

	return &Pagination{Base: base, Page: page, Start: start, BarSize: barSize, TotalPages: totalPages}
}
