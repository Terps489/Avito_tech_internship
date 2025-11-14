package postgres

import (
	"database/sql"

	"github.com/terps489/avito_tech_internship/internal/domain"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) GetByID(id domain.TeamID) (*domain.Team, error) {
	const query = `
		SELECT id, name
		FROM teams
		WHERE id = $1
	`

	var t domain.Team
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
