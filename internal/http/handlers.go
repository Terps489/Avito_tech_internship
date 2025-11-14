package http

import (
	"net/http"
)

// For testing purposes only
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// ---------- Teams ----------

func (s *Server) handleTeamAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

func (s *Server) handleTeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	// query?team_name=...

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

// ---------- Users ----------

func (s *Server) handleUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

func (s *Server) handleUserGetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

// ---------- Pull Requests ----------

func (s *Server) handlePullRequestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

func (s *Server) handlePullRequestMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}

func (s *Server) handlePullRequestReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrorCodeNotFound,
			Message: "not implemented yet",
		},
	})
}
