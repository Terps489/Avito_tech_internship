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

func (r *TeamRepository) GetByName(name domain.TeamName) (*domain.Team, error) {
	const query = `
		SELECT team_name
		FROM teams
		WHERE team_name = $1
	`

	var t domain.Team
	err := r.db.QueryRow(query, name).Scan(&t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
