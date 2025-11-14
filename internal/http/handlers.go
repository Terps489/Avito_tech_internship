package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "invalid json body",
			},
		})
		return
	}

	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "user_id is required",
			},
		})
		return
	}

	u, err := s.service.SetUserIsActive(domain.UserID(req.UserID), req.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotFound,
					Message: "user not found",
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
		User UserDTO `json:"user"`
	}{
		User: toUserDTO(u),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleUserGetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "user_id query param is required",
			},
		})
		return
	}

	prs, err := s.service.ListPullRequestsForReviewer(domain.UserID(userID))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "internal error: " + err.Error(),
			},
		})
		return
	}

	resp := struct {
		UserID       string                `json:"user_id"`
		PullRequests []PullRequestShortDTO `json:"pull_requests"`
	}{
		UserID:       userID,
		PullRequests: make([]PullRequestShortDTO, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, PullRequestShortDTO{
			ID:       string(pr.ID),
			Name:     pr.Title,
			AuthorID: string(pr.AuthorID),
			Status:   string(pr.Status),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

// ---------- Pull Requests ----------

func (s *Server) handlePullRequestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "invalid json body",
			},
		})
		return
	}

	if req.ID == "" || req.Name == "" || req.Author == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "pull_request_id, pull_request_name and author_id are required",
			},
		})
		return
	}

	pr, err := s.service.CreatePullRequestWithID(
		domain.PullRequestID(req.ID),
		req.Name,
		domain.UserID(req.Author),
	)
	if err != nil {
		if errors.Is(err, app.ErrPRExists) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodePRExists,
					Message: "PR id already exists",
				},
			})
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotFound,
					Message: "author not found",
				},
			})
			return
		}

		if errors.Is(err, app.ErrAuthorNotActive) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNoCandidate,
					Message: "author is not active",
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
		PR PullRequestDTO `json:"pr"`
	}{
		PR: toPullRequestDTO(pr),
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handlePullRequestMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "invalid json body",
			},
		})
		return
	}

	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "pull_request_id is required",
			},
		})
		return
	}

	pr, err := s.service.MergePullRequest(domain.PullRequestID(req.ID))
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotFound,
					Message: "pull request not found",
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
		PR PullRequestDTO `json:"pr"`
	}{
		PR: toPullRequestDTO(pr),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handlePullRequestReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	var req ReassignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "invalid json body",
			},
		})
		return
	}

	if req.PRID == "" || req.OldUserID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorPayload{
				Code:    ErrorCodeNotFound,
				Message: "pull_request_id and old_user_id are required",
			},
		})
		return
	}

	pr, replacedBy, err := s.service.ReassignReviewer(
		domain.PullRequestID(req.PRID),
		domain.UserID(req.OldUserID),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotFound,
					Message: "pull request or user not found",
				},
			})
			return
		}

		if errors.Is(err, app.ErrPRAlreadyMerged) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodePRMerged,
					Message: "cannot reassign on merged PR",
				},
			})
			return
		}

		if errors.Is(err, app.ErrReviewerNotAssigned) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNotAssigned,
					Message: "reviewer is not assigned to this PR",
				},
			})
			return
		}

		if errors.Is(err, app.ErrNoAvailableReviewers) {
			writeJSON(w, http.StatusConflict, ErrorResponse{
				Error: ErrorPayload{
					Code:    ErrorCodeNoCandidate,
					Message: "no active replacement candidate in team",
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
		PR         PullRequestDTO `json:"pr"`
		ReplacedBy string         `json:"replaced_by"`
	}{
		PR:         toPullRequestDTO(pr),
		ReplacedBy: string(replacedBy),
	}

	writeJSON(w, http.StatusOK, resp)
}

func toPullRequestDTO(pr *domain.PullRequest) PullRequestDTO {
	dto := PullRequestDTO{
		ID:                string(pr.ID),
		Name:              pr.Title,
		AuthorID:          string(pr.AuthorID),
		Status:            string(pr.Status),
		AssignedReviewers: make([]string, 0, len(pr.ReviewerIDs)),
	}

	for _, id := range pr.ReviewerIDs {
		dto.AssignedReviewers = append(dto.AssignedReviewers, string(id))
	}

	if !pr.CreatedAt.IsZero() {
		s := pr.CreatedAt.UTC().Format(time.RFC3339)
		dto.CreatedAt = &s
	}

	if pr.MergedAt != nil {
		s := pr.MergedAt.UTC().Format(time.RFC3339)
		dto.MergedAt = &s
	}

	return dto
}

func toUserDTO(u *domain.User) UserDTO {
	return UserDTO{
		UserID:   string(u.ID),
		Username: u.Username,
		TeamName: string(u.TeamName),
		IsActive: u.IsActive,
	}
}
