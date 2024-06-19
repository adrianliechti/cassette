package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/cors"

	"cassette/pkg/repository"
	_ "embed"
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
	handler http.Handler

	*repository.Repository

	last string
}

func New(r *repository.Repository) *Server {
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
		handler: handler,

		Repository: r,
	}

	mux.HandleFunc("POST /events", s.handleEvents)
	mux.HandleFunc("GET /cassette.min.cjs", s.handleScript)

	mux.HandleFunc("GET /sessions", s.handleSessions)
	mux.HandleFunc("GET /sessions/{session}", s.handleSession)
	mux.HandleFunc("GET /sessions/{session}/events", s.handleSessionEvents)

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

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
		Events []repository.Event `json:"events"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionID := s.getSessionID(r)
	s.last = sessionID

	if err := s.AppendSessionEvents(sessionID, body.Events...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: sessionID,

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

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if id == "default" && s.last != "" {
		id = s.last
	}

	session, err := s.Session(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	result := Session{
		ID:      session.ID,
		Created: session.Created,
	}

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleSessionEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if id == "default" && s.last != "" {
		id = s.last
	}

	events, err := s.SessionEvents(id)

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

	return uuid.New().String()
}
