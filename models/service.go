package models

type Service struct {
	ID      int    `json:"id";db:"id"`
	Name    string `json:"name";db:"name"`
	PID     int    `json:"pid";db:"pid"`
	Status  int    `json:"status";db:"status"`
	Deleted int    `json:"deleted";db:"deleted"`
}
