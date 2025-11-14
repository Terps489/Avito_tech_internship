package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/terps489/avito_tech_internship/internal/app"
	"github.com/terps489/avito_tech_internship/internal/domain"
)

// --- Health ---

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

	var body TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "invalid json body",
			},
		})
		return
	}

	if body.TeamName == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "team_name is required",
			},
		})
		return
	}

	var members []domain.User
	for _, m := range body.Members {
		members = append(members, domain.User{
			ID:       domain.UserID(m.UserID),
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: domain.TeamName(body.TeamName),
		})
	}

	team, membersFromDB, err := s.service.CreateTeamWithMembers(domain.TeamName(body.TeamName), members)
	if err != nil {
		if errors.Is(err, app.ErrTeamExists) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeTeamExists,
					Message: "team_name already exists",
				},
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "internal error: " + err.Error(),
			},
		})
		return
	}

	resp := struct {
		Team TeamDTO `json:"team"`
	}{
		Team: TeamDTO{
			TeamName: string(team.Name),
			Members:  make([]TeamMemberDTO, 0, len(membersFromDB)),
		},
	}

	for _, m := range membersFromDB {
		resp.Team.Members = append(resp.Team.Members, TeamMemberDTO{
			UserID:   string(m.ID),
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handleTeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "team_name query param is required",
			},
		})
		return
	}

	team, members, err := s.service.GetTeamWithMembers(domain.TeamName(teamName))
	if err != nil {
		if errors.Is(err, app.ErrTeamNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotFound,
					Message: "team not found",
				},
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "internal error: " + err.Error(),
			},
		})
		return
	}

	resp := TeamDTO{
		TeamName: string(team.Name),
		Members:  make([]TeamMemberDTO, 0, len(members)),
	}

	for _, m := range members {
		resp.Members = append(resp.Members, TeamMemberDTO{
			UserID:   string(m.ID),
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	writeJSON(w, http.StatusOK, resp)
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
