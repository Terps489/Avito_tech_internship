package app

import (
	"errors"
	"math/rand"
	"time"

	"github.com/terps489/avito_tech_internship/internal/domain"
)

var (
	ErrPRAlreadyMerged      = errors.New("pull request is already merged")
	ErrReviewerNotAssigned  = errors.New("reviewer is not assigned to pull request")
	ErrNoAvailableReviewers = errors.New("no available reviewers in team")
	ErrAuthorNotActive      = errors.New("author is not active")
)

type UserRepository interface {
	GetByID(id domain.UserID) (*domain.User, error)
	ListActiveByTeam(teamName domain.TeamName) ([]domain.User, error)
}

type TeamRepository interface {
	GetByName(name domain.TeamName) (*domain.Team, error)
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
	rnd   *rand.Rand
}

func NewService(u UserRepository, t TeamRepository, p PullRequestRepository) *Service {
	return &Service{
		users: u,
		teams: t,
		prs:   p,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Service) CreatePullRequest(authorID domain.UserID, title string) (*domain.PullRequest, error) {
	author, err := s.users.GetByID(authorID)
	if err != nil {
		return nil, err
	}

	if !author.IsActive {
		return nil, ErrAuthorNotActive
	}

	candidates, err := s.users.ListActiveByTeam(author.TeamName)
	if err != nil {
		return nil, err
	}

	var reviewerPool []domain.UserID
	for _, u := range candidates {
		if u.ID == author.ID {
			continue
		}
		reviewerPool = append(reviewerPool, u.ID)
	}

	s.shuffleUserIDs(reviewerPool)

	var reviewers []domain.UserID
	for i := 0; i < len(reviewerPool) && len(reviewers) < 2; i++ {
		reviewers = append(reviewers, reviewerPool[i])
	}

	pr := &domain.PullRequest{
		Title:       title,
		AuthorID:    authorID,
		Status:      domain.PRStatusOpen,
		ReviewerIDs: reviewers,
	}

	if err := s.prs.Create(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) ReassignReviewer(prID domain.PullRequestID, oldReviewerID domain.UserID) (*domain.PullRequest, error) {
	pr, err := s.prs.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, ErrPRAlreadyMerged
	}

	idx := -1
	for i, id := range pr.ReviewerIDs {
		if id == oldReviewerID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, ErrReviewerNotAssigned
	}

	reviewer, err := s.users.GetByID(oldReviewerID)
	if err != nil {
		return nil, err
	}

	candidates, err := s.users.ListActiveByTeam(reviewer.TeamName)
	if err != nil {
		return nil, err
	}

	exclude := make(map[domain.UserID]struct{}, len(pr.ReviewerIDs)+1)
	for _, id := range pr.ReviewerIDs {
		exclude[id] = struct{}{}
	}
	exclude[oldReviewerID] = struct{}{}

	var pool []domain.UserID
	for _, u := range candidates {
		if _, banned := exclude[u.ID]; banned {
			continue
		}
		pool = append(pool, u.ID)
	}

	if len(pool) == 0 {
		return nil, ErrNoAvailableReviewers
	}

	newReviewerID := pool[s.rnd.Intn(len(pool))]

	pr.ReviewerIDs[idx] = newReviewerID

	if err := s.prs.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) MergePullRequest(prID domain.PullRequestID) (*domain.PullRequest, error) {
	pr, err := s.prs.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	pr.Status = domain.PRStatusMerged

	if err := s.prs.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) shuffleUserIDs(ids []domain.UserID) {
	if len(ids) < 2 {
		return
	}
	s.rnd.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
}
