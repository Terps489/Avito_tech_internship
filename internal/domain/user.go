package domain

type UserID string
type TeamName string

type User struct {
	ID       UserID
	Username string
	IsActive bool
	TeamName TeamName
}
