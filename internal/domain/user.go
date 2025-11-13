package domain

type UserID int64

type User struct {
	ID       UserID
	Name     string
	IsActive bool
	TeamID   TeamID
}
