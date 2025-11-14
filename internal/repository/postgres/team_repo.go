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

func (r *TeamRepository) Create(name domain.TeamName) error {
	const query = `
		INSERT INTO teams (team_name)
		VALUES ($1)
	`
	_, err := r.db.Exec(query, name)
	return err
}

func (r *TeamRepository) Exists(name domain.TeamName) (bool, error) {
	const query = `
		SELECT 1
		FROM teams
		WHERE team_name = $1
	`
	var dummy int
	err := r.db.QueryRow(query, name).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *TeamRepository) ListMembers(name domain.TeamName) ([]domain.User, error) {
	const query = `
		SELECT user_id, username, is_active, team_name
		FROM users
		WHERE team_name = $1
		ORDER BY user_id
	`

	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
