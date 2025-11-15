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
		SELECT user_id, username, is_active, team_name
		FROM users
		WHERE user_id = $1
	`

	var u domain.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) ListActiveByTeam(teamName domain.TeamName) ([]domain.User, error) {
	const query = `
		SELECT user_id, username, is_active, team_name
		FROM users
		WHERE team_name = $1 AND is_active = TRUE
	`

	rows, err := r.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

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

func (r *UserRepository) UpsertUsersForTeam(teamName domain.TeamName, users []domain.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const query = `
		INSERT INTO users (user_id, username, is_active, team_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET username = EXCLUDED.username,
		    is_active = EXCLUDED.is_active,
		    team_name = EXCLUDED.team_name
	`

	for _, u := range users {
		if _, err := tx.Exec(query, u.ID, u.Username, u.IsActive, teamName); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *UserRepository) SetIsActive(id domain.UserID, active bool) error {
	const query = `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
	`

	res, err := r.db.Exec(query, id, active)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
