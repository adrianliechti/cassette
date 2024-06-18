package server

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/cors"

	"cassette/pkg/repository"
)

const (
	cookieName = "casette-session"
)

var (
	//go:embed cassette.min.cjs
	cassetteRecord string

	//go:embed cassette.js
	cassetteClient string
)

type Server struct {
	handler http.Handler

	*repository.Repository
}

func New(r *repository.Repository) *Server {
	mux := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := cors.Handler(mux)

	s := &Server{
		handler: handler,

		Repository: r,
	}

	mux.HandleFunc("GET /", s.handleIndex)

	mux.HandleFunc("POST /events", s.handleEvents)

	mux.HandleFunc("GET /sessions", s.handleSessions)
	mux.HandleFunc("GET /sessions/{session}/events", s.handleSessionEvents)

	mux.HandleFunc("GET /cassette.min.cjs", s.handleScript)

	return s
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) handleScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")

	var result strings.Builder

	result.WriteString(cassetteRecord)
	result.WriteString("\n")
	result.WriteString(cassetteClient)

	w.Write([]byte(result.String()))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Events []Event `json:"events"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := s.Session(s.getSessionID(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.AppendEvents(body.Events); err != nil {
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
	sessions, err := s.Sessions()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	result := make([]Session, 0)

	for _, s := range sessions {
		result = append(result, Session{
			ID: s.ID,

			Created: s.Created,
		})
	}

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleSessionEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	session, err := s.Session(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	events := session.Events()

	json.NewEncoder(w).Encode(events)
}

func (s *Server) getSessionID(r *http.Request) string {
	cookie, _ := r.Cookie(cookieName)

	if cookie != nil && cookie.Value != "" {
		return cookie.Value
	}

	return uuid.New().String()
}
