package postgres

import (
	"database/sql"
	"time"

	"github.com/terps489/avito_tech_internship/internal/domain"
)

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Create(pr *domain.PullRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const insertPR = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := tx.Exec(insertPR, pr.ID, pr.Title, pr.AuthorID, pr.Status); err != nil {
		return err
	}

	if len(pr.ReviewerIDs) > 0 {
		const insertReviewer = `
			INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
			VALUES ($1, $2)
		`
		for _, reviewerID := range pr.ReviewerIDs {
			if _, err := tx.Exec(insertReviewer, pr.ID, reviewerID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PullRequestRepository) GetByID(id domain.PullRequestID) (*domain.PullRequest, error) {
	const queryPR = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr domain.PullRequest
	var mergedAt sql.NullTime

	if err := r.db.QueryRow(queryPR, id).
		Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt); err != nil {
		return nil, err
	}

	if mergedAt.Valid {
		t := mergedAt.Time
		pr.MergedAt = &t
	}

	const queryReviewers = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pr_id = $1
	`

	rows, err := r.db.Query(queryReviewers, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var rid domain.UserID
		if err := rows.Scan(&rid); err != nil {
			return nil, err
		}
		pr.ReviewerIDs = append(pr.ReviewerIDs, rid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) Update(pr *domain.PullRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const updatePR = `
		UPDATE pull_requests
		SET pull_request_name = $1,
		    author_id = $2,
		    status = $3,
		    merged_at = $4
		WHERE pull_request_id = $5
	`

	var mergedAt interface{}
	if pr.Status == domain.PRStatusMerged {
		if pr.MergedAt == nil {
			now := time.Now().UTC()
			pr.MergedAt = &now
		}
		mergedAt = pr.MergedAt
	} else {
		mergedAt = nil
	}

	if _, err := tx.Exec(updatePR, pr.Title, pr.AuthorID, pr.Status, mergedAt, pr.ID); err != nil {
		return err
	}

	const deleteReviewers = `
		DELETE FROM pull_request_reviewers
		WHERE pr_id = $1
	`
	if _, err := tx.Exec(deleteReviewers, pr.ID); err != nil {
		return err
	}

	const insertReviewer = `
		INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2)
	`
	for _, reviewerID := range pr.ReviewerIDs {
		if _, err := tx.Exec(insertReviewer, pr.ID, reviewerID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PullRequestRepository) Exists(id domain.PullRequestID) (bool, error) {
	const query = `
		SELECT 1
		FROM pull_requests
		WHERE pull_request_id = $1
	`
	var dummy int
	err := r.db.QueryRow(query, id).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PullRequestRepository) ListByReviewer(userID domain.UserID) ([]domain.PullRequest, error) {
	const query = `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pull_request_reviewers r ON r.pr_id = pr.pull_request_id
		WHERE r.reviewer_id = $1
		ORDER BY pr.pull_request_id
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PullRequestRepository) GetReviewerAssignmentStats() ([]domain.ReviewerAssignmentStat, error) {
	const query = `
		SELECT reviewer_id, COUNT(*) as cnt
		FROM pull_request_reviewers
		GROUP BY reviewer_id
		ORDER BY reviewer_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var stats []domain.ReviewerAssignmentStat
	for rows.Next() {
		var s domain.ReviewerAssignmentStat
		if err := rows.Scan(&s.UserID, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
