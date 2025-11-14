package postgres

import (
	"database/sql"

	"github.com/terps489/avito_tech_internship/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id domain.UserID) (*domain.User, error) {
	const query = `
		SELECT id, name, is_active, team_id
		FROM users
		WHERE id = $1
	`

	var u domain.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) ListActiveByTeam(teamID domain.TeamID) ([]domain.User, error) {
	const query = `
		SELECT id, name, is_active, team_id
		FROM users
		WHERE team_id = $1 AND is_active = TRUE
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
