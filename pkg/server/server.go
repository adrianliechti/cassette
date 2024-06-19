package server

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/rs/cors"

	"cassette/config"
	"cassette/pkg/repository"
	"cassette/pkg/storage"
)

const (
	cookieName = "casette-session"
)

var (
	//go:embed cassette.js
	jsCassette string

	//go:embed record.umd.min.cjs
	jsRecord string
)

type Server struct {
	*config.Config

	handler http.Handler

	lastSession string
}

func New(config *config.Config) *Server {
	mux := http.NewServeMux()

	cors.Default()

	cors := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := cors.Handler(mux)

	s := &Server{
		Config: config,

		handler: handler,
	}

	mux.HandleFunc("POST /events", s.handleEvents)
	mux.HandleFunc("GET /cassette.min.cjs", s.handleScript)

	mux.HandleFunc("GET /sessions", s.handleSessions)
	mux.HandleFunc("GET /sessions/{session}", s.handleSession)
	mux.HandleFunc("GET /sessions/{session}/events", s.handleSessionEvents)

	mux.Handle("/", http.FileServer(http.Dir("./public")))

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")

	var result bytes.Buffer

	result.WriteString(jsRecord)
	result.WriteString("\n")
	result.WriteString(jsCassette)

	w.Write(result.Bytes())
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Events []storage.Event `json:"events"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	var session *repository.Session

	if id := s.getSessionID(r); id != "" {
		session, err = s.Repository.Session(id)
	}

	if session == nil {
		session, err = s.Repository.CreateSession()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.lastSession = session.ID

	if err := s.Storage.AppendEvents(session.ID, body.Events...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: session.ID,

		Path: "/",

		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := s.Repository.Sessions()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	json.NewEncoder(w).Encode(sessions)
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if id == "default" && s.lastSession != "" {
		id = s.lastSession
	}

	session, err := s.Repository.Session(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func (s *Server) handleSessionEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if id == "default" && s.lastSession != "" {
		id = s.lastSession
	}

	session, err := s.Repository.Session(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	events, err := s.Storage.Events(session.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (s *Server) getSessionID(r *http.Request) string {
	cookie, _ := r.Cookie(cookieName)

	if cookie != nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}
