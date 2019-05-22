package models

type User struct {
	ID       int    `json:"id";db:"id"`
	Username string `json:"username";db:"username"`
	Role     int    `json:"role";db:"role"`
	Deleted  int    `json:"deleted";db:"deleted"`
}
