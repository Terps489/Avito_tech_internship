package domain

import "time"

type PullRequestID string

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID          PullRequestID
	Title       string
	AuthorID    UserID
	Status      PRStatus
	ReviewerIDs []UserID
	CreatedAt   time.Time
	MergedAt    *time.Time
}
