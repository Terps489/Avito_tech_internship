package app

import "github.com/terps489/avito_tech_internship/internal/domain"

type UserRepository interface {
	GetByID(id domain.UserID) (*domain.User, error)
	ListActiveByTeam(teamID domain.TeamID) ([]domain.User, error)
}

type TeamRepository interface {
	GetByID(id domain.TeamID) (*domain.Team, error)
}

type PullRequestRepository interface {
	Create(pr *domain.PullRequest) error
	GetByID(id domain.PullRequestID) (*domain.PullRequest, error)
	Update(pr *domain.PullRequest) error
}

type Service struct {
	users UserRepository
	teams TeamRepository
	prs   PullRequestRepository
}

func NewService(u UserRepository, t TeamRepository, p PullRequestRepository) *Service {
	return &Service{
		users: u,
		teams: t,
		prs:   p,
	}
}

// Создание Pull Request
func (s *Service) CreatePullRequest(authorID domain.UserID, title string) (*domain.PullRequest, error) {
	//
	return nil, nil
}

// Переназначение
func (s *Service) ReassignReviewer(prID domain.PullRequestID, oldReviewerID domain.UserID) (*domain.PullRequest, error) {
	//
	return nil, nil
}

//Merge Pull Request
func (s *Service) MergePullRequest(prID domain.PullRequestID) (*domain.PullRequest, error) {
	//
	return nil, nil
}
