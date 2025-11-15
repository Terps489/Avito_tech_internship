package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/terps489/avito_tech_internship/internal/app"
)

type Server struct {
	addr    string
	service *app.Service
	mux     *http.ServeMux
}

func NewServer(addr string, svc *app.Service) *Server {
	s := &Server{
		addr:    addr,
		service: svc,
		mux:     http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Run() error {
	log.Printf("starting http server on %s", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)

	// Teams
	s.mux.HandleFunc("/team/add", s.handleTeamAdd)
	s.mux.HandleFunc("/team/get", s.handleTeamGet)

	// Stats
	s.mux.HandleFunc("/stats/assignments", s.handleStatsAssignments)

	// Users
	s.mux.HandleFunc("/users/setIsActive", s.handleUserSetIsActive)
	s.mux.HandleFunc("/users/getReview", s.handleUserGetReview)

	// Pull Requests
	s.mux.HandleFunc("/pullRequest/create", s.handlePullRequestCreate)
	s.mux.HandleFunc("/pullRequest/merge", s.handlePullRequestMerge)
	s.mux.HandleFunc("/pullRequest/reassign", s.handlePullRequestReassign)
}

// ---------- Helpers ----------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
		"error": map[string]any{
			"code":    "METHOD_NOT_ALLOWED",
			"message": "method not allowed",
		},
	})
}
