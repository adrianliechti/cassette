package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

func main() {
	s := New()

	log.Fatal(http.ListenAndServe(":8080", s))
}

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

	sessions map[string]*Session

	lastSessionID string
}

type Session struct {
	ID string `json:"id,omitempty"`

	Timestamp time.Time `json:"timestamp,omitempty"`

	Events []any `json:"events,omitempty"`
}

func New() *Server {
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

		sessions: make(map[string]*Session),
	}

	mux.HandleFunc("GET /", s.handleIndex)

	mux.HandleFunc("POST /events", s.handleEvents)

	mux.HandleFunc("GET /sessions", s.handleSessions)
	mux.HandleFunc("GET /sessions/{session}", s.handleSession)

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

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	var events struct {
		Events []any `json:"events"`
	}

	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := s.getSession(r)

	s.lastSessionID = session.ID

	session.Events = append(session.Events, events.Events...)

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: session.ID,

		Path:     "/",
		HttpOnly: true,
	})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	result := make([]Session, 0)

	for _, session := range s.sessions {
		result = append(result, Session{
			ID: session.ID,

			Timestamp: session.Timestamp,
		})
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if id == "default" {
		id = s.lastSessionID
	}

	session, ok := s.sessions[id]

	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) getSession(r *http.Request) *Session {
	id := uuid.New().String()

	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		id = cookie.Value
	}

	if session, ok := s.sessions[id]; ok {
		return session
	}

	session := &Session{
		ID:        id,
		Timestamp: time.Now(),

		Events: make([]any, 0),
	}

	s.sessions[id] = session

	return session
}
