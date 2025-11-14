package postgres

import (
	"database/sql"

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
	defer tx.Rollback()

	const insertPR = `
		INSERT INTO pull_requests (title, author_id, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	if err := tx.QueryRow(insertPR, pr.Title, pr.AuthorID, pr.Status).Scan(&pr.ID); err != nil {
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
		SELECT id, title, author_id, status
		FROM pull_requests
		WHERE id = $1
	`
	var pr domain.PullRequest
	if err := r.db.QueryRow(queryPR, id).
		Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
		return nil, err
	}

	// reviewers
	const queryReviewers = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pr_id = $1
	`

	rows, err := r.db.Query(queryReviewers, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	defer tx.Rollback()

	const updatePR = `
		UPDATE pull_requests
		SET title = $1, author_id = $2, status = $3
		WHERE id = $4
	`
	if _, err := tx.Exec(updatePR, pr.Title, pr.AuthorID, pr.Status, pr.ID); err != nil {
		return err
	}

	// delete reviewers and re-insert
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
