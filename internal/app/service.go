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
	ErrTeamExists           = errors.New("team already exists")
	ErrTeamNotFound         = errors.New("team not found")
	ErrPRExists             = errors.New("pull request already exists")
)

// ---------- Репозитории ----------

type UserRepository interface {
	GetByID(id domain.UserID) (*domain.User, error)
	ListActiveByTeam(teamName domain.TeamName) ([]domain.User, error)
	UpsertUsersForTeam(teamName domain.TeamName, users []domain.User) error
}

type TeamRepository interface {
	GetByName(name domain.TeamName) (*domain.Team, error)
	Create(name domain.TeamName) error
	Exists(name domain.TeamName) (bool, error)
	ListMembers(name domain.TeamName) ([]domain.User, error)
}

type PullRequestRepository interface {
	Create(pr *domain.PullRequest) error
	GetByID(id domain.PullRequestID) (*domain.PullRequest, error)
	Update(pr *domain.PullRequest) error
	Exists(id domain.PullRequestID) (bool, error)
}

// ---------- Service ----------

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

// ---------- Команды ----------

func (s *Service) CreateTeamWithMembers(teamName domain.TeamName, members []domain.User) (*domain.Team, []domain.User, error) {
	exists, err := s.teams.Exists(teamName)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrTeamExists
	}

	if err := s.teams.Create(teamName); err != nil {
		return nil, nil, err
	}

	for i := range members {
		members[i].TeamName = teamName
	}

	if err := s.users.UpsertUsersForTeam(teamName, members); err != nil {
		return nil, nil, err
	}

	team, err := s.teams.GetByName(teamName)
	if err != nil {
		return nil, nil, err
	}

	membersFromDB, err := s.teams.ListMembers(teamName)
	if err != nil {
		return nil, nil, err
	}

	return team, membersFromDB, nil
}

func (s *Service) CreatePullRequestWithID(
	id domain.PullRequestID,
	name string,
	authorID domain.UserID,
) (*domain.PullRequest, error) {

	exists, err := s.prs.Exists(id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPRExists
	}

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

	reviewers := make([]domain.UserID, 0, 2)
	for i := 0; i < len(reviewerPool) && len(reviewers) < 2; i++ {
		reviewers = append(reviewers, reviewerPool[i])
	}

	pr := &domain.PullRequest{
		ID:          id,
		Title:       name,
		AuthorID:    authorID,
		Status:      domain.PRStatusOpen,
		ReviewerIDs: reviewers,
	}

	if err := s.prs.Create(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) GetTeamWithMembers(teamName domain.TeamName) (*domain.Team, []domain.User, error) {
	exists, err := s.teams.Exists(teamName)
	if err != nil {
		return nil, nil, err
	}
	if !exists {
		return nil, nil, ErrTeamNotFound
	}

	team, err := s.teams.GetByName(teamName)
	if err != nil {
		return nil, nil, err
	}

	members, err := s.teams.ListMembers(teamName)
	if err != nil {
		return nil, nil, err
	}

	return team, members, nil
}

// ---------- PR: создание / переназначение / merge ----------

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
