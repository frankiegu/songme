package web

import (
	"net/http"

	"github.com/emredir/songme/models"
	"github.com/gorilla/mux"
)

// NewServer returns a new server.
func NewServer(userStore models.UserStore, songStore models.SongStore) *Server {
	server := &Server{
		middleware: &Middleware{userStore},
		auth:       &AuthHandler{userStore},
		admin:      &AdminHandler{songStore},
		song:       &SongHandler{songStore},
		router:     mux.NewRouter(),
	}
	server.buildRoutes()
	return server
}

// Server is our HTTP server with routes for all our endpoints.
type Server struct {
	middleware *Middleware
	auth       *AuthHandler
	admin      *AdminHandler
	song       *SongHandler
	router     *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) buildRoutes() {
	// Routes
	// FileServer
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Main router
	s.router.HandleFunc("/login", s.auth.RenderLogin).Methods("GET")
	s.router.HandleFunc("/login", s.auth.Login).Methods("POST")
	s.router.HandleFunc("/logout", s.auth.Logout).Methods("POST")
	s.router.HandleFunc("/", s.song.Index).Methods("GET")
	s.router.HandleFunc("/add", s.song.New).Methods("GET")
	s.router.HandleFunc("/add", s.song.Create).Methods("POST")

	// Admin router
	adminRouter := s.router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/dashboard", s.admin.Dashboard).Methods("GET")

	// Songs router
	songsRouter := s.router.PathPrefix("/songs").Subrouter()
	songsRouter.HandleFunc("/{id}", s.song.Confirm).Methods("PUT")
	songsRouter.HandleFunc("/candidate", s.song.Candidates).Methods("GET")
	songsRouter.HandleFunc("/production", s.song.Productions).Methods("GET")
	songsRouter.HandleFunc("/candidate/{id}", s.song.DeleteCandidate).Methods("DELETE")
	songsRouter.HandleFunc("/production/{id}", s.song.DeleteProduction).Methods("DELETE")

	// Recover panics
	s.router.Use(s.middleware.PanicRecovery)

	// Authorize admin router
	adminRouter.Use(s.middleware.Authorize)

	// Authorize songs router
	songsRouter.Use(s.middleware.Authorize)
}
