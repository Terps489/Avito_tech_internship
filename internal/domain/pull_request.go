package domain

type PullRequestID int64

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
	ReviewerIDs []UserID // max reviewers: 2
}
