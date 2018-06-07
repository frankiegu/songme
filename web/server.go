package web

import (
	"net/http"

	"github.com/emredir/songme"
	"github.com/emredir/songme/models"

	"github.com/gorilla/mux"
)

// Store packs store interfaces.
type Store struct {
	User models.UserStore
	Role models.RoleStore
	Song models.SongStore
}

// Interactor packs interactor interfaces.
type Interactor struct {
	Auth AuthInteractor
	Song SongInteractor
}

// NewServer returns a new server.
func NewServer(store Store, interactor Interactor) *Server {
	middleware := Middleware{
		userStore: store.User,
	}
	auth := AuthHandler{
		AuthInteractor: interactor.Auth,
		UsernameLength: songme.GetConfig().UsernameLength,
		PasswordLength: songme.GetConfig().PasswordLength,
	}
	main := MainHandler{
		songStore: store.Song,
	}
	song := SongHandler{
		songInteractor:    interactor.Song,
		descriptionLength: songme.GetConfig().SongDescriptionLength,
		songsPerPage:      songme.GetConfig().SongsPerPage,
	}
	user := UserHandler{
		userStore:    store.User,
		songStore:    store.Song,
		songsPerPage: songme.GetConfig().SongsPerPage,
	}
	admin := AdminHandler{
		userStore:    store.User,
		songStore:    store.Song,
		songsPerPage: songme.GetConfig().SongsPerPage,
		usersPerPage: songme.GetConfig().UsersPerPage,
	}

	server := &Server{
		middleware: &middleware,
		auth:       &auth,
		main:       &main,
		song:       &song,
		user:       &user,
		admin:      &admin,
		router:     mux.NewRouter(),
	}
	server.buildRoutes()

	return server
}

// Server is our HTTP server with routes for all our endpoints.
type Server struct {
	middleware *Middleware
	auth       *AuthHandler
	main       *MainHandler
	song       *SongHandler
	user       *UserHandler
	admin      *AdminHandler
	router     *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) buildRoutes() {
	// Routes
	// FileServer
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Main routes
	s.router.HandleFunc("/", s.main.Home).Methods("GET")

	// Auth routes
	s.router.HandleFunc("/signup", s.auth.RenderSignup).Methods("GET")
	s.router.HandleFunc("/signup", s.auth.Signup).Methods("POST")
	s.router.HandleFunc("/signin", s.auth.RenderSignin).Methods("GET")
	s.router.HandleFunc("/signin", s.auth.Signin).Methods("POST")
	// Authorized auth routes
	s.router.Handle("/logout", s.middleware.Authorize(http.HandlerFunc(s.auth.Logout))).Methods("GET", "POST")
	s.router.Handle("/update-password", s.middleware.Authorize(http.HandlerFunc(s.auth.RenderUpdatePassword))).Methods("GET")
	s.router.Handle("/update-password", s.middleware.Authorize(http.HandlerFunc(s.auth.UpdatePassword))).Methods("POST")
	s.router.Handle("/update-email", s.middleware.Authorize(http.HandlerFunc(s.auth.RenderUpdateEmail))).Methods("GET")
	s.router.Handle("/update-email", s.middleware.Authorize(http.HandlerFunc(s.auth.UpdateEmail))).Methods("POST")

	// Song router
	s.router.HandleFunc("/recommend", s.song.New).Methods("GET")
	s.router.HandleFunc("/recommend", s.song.Create).Methods("POST")
	s.router.HandleFunc("/songs", s.song.Songs).Methods("GET")
	s.router.HandleFunc("/songs/page/{page:[0-9]+}", s.song.Songs).Methods("GET")

	// User router
	s.router.HandleFunc("/user/{username}", s.user.Profile).Methods("GET")
	s.router.HandleFunc("/user/{username}/page/{page:[0-9]+}", s.user.Profile).Methods("GET")

	// Admin router
	adminRouter := s.router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/dashboard", s.admin.Dashboard).Methods("GET")
	adminRouter.HandleFunc("/dashboard/page/{page:[0-9]+}", s.admin.Dashboard).Methods("GET")
	adminRouter.HandleFunc("/dashboard/productions", s.admin.Productions).Methods("GET")
	adminRouter.HandleFunc("/dashboard/productions/{page:[0-9]+}", s.admin.Productions).Methods("GET")
	adminRouter.HandleFunc("/dashboard/users", s.admin.Users).Methods("GET")
	adminRouter.HandleFunc("/dashboard/users/{page:[0-9]+}", s.admin.Users).Methods("GET")
	adminRouter.HandleFunc("/confirm/{id}", s.admin.ConfirmSong).Methods("POST")
	adminRouter.HandleFunc("/delete/{id}", s.admin.DeleteSong).Methods("POST")

	// Recover panics
	s.router.Use(s.middleware.PanicRecovery)

	// User via session
	s.router.Use(s.middleware.UserViaSession)

	// Authorize admin router
	adminRouter.Use(s.middleware.Admin)
}
