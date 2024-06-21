package server

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"os"

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

	handler    http.Handler
	filesystem fs.FS
}

func New(config *config.Config) *Server {
	mux := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	s := &Server{
		Config: config,

		handler:    cors.Handler(mux),
		filesystem: os.DirFS("./public"),
	}

	mux.HandleFunc("POST /events", s.handleEvents)
	mux.HandleFunc("GET /cassette.min.cjs", s.handleScript)

	mux.HandleFunc("GET /sessions", s.handleAuth(s.handleSessions))
	mux.HandleFunc("GET /sessions/{session}", s.handleAuth(s.handleSession))
	mux.HandleFunc("GET /sessions/{session}/events", s.handleAuth(s.handleSessionEvents))
	mux.HandleFunc("DELETE /sessions/{session}", s.handleAuth(s.handleSessionDelete))

	mux.HandleFunc("/", s.handleAuth(s.handleUI))

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Username == "" || s.Password == "" {
			next(w, r)
			return
		}

		if username, password, ok := r.BasicAuth(); ok {
			if username == s.Username && password == s.Password {
				next(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Cassette - Admin UI", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/help" || r.URL.Path == "/help/" {
		http.ServeFileFS(w, r, s.filesystem, "index.html")
		return
	}

	handler := http.FileServerFS(s.filesystem)
	handler.ServeHTTP(w, r)
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

	if id := getSessionID(r); id != "" {
		session, err = s.Repository.Session(id)
	}

	if session == nil {
		info := &repository.SessionInfo{
			Origin:  getOrigin(r),
			Address: getAddress(r),

			UserAgent: r.UserAgent(),
		}

		session, err = s.Repository.CreateSession(info)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Storage.AppendEvents(session.ID, body.Events...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setSessionID(w, r, session.ID)
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

	session, err := s.Repository.Session(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func (s *Server) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

	if err := s.Storage.DeleteDelete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Repository.DeleteSession(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSessionEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("session")

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

func getOrigin(r *http.Request) string {
	if val := r.Header.Get("Origin"); val != "" {
		return val
	}

	if val := r.Header.Get("Referer"); val != "" {
		return val
	}

	return ""
}

func getAddress(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func getSessionID(r *http.Request) string {
	cookie, _ := r.Cookie(cookieName)

	if cookie != nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func setSessionID(w http.ResponseWriter, r *http.Request, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: id,

		Path: "/",

		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})
}
