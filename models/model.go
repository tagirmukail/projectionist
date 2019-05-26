package models

import "database/sql"

const (
	NotDeleted int = iota
	Deleted
)

type Model interface {
	Validate() error
	IsExist(*sql.DB) (error, bool)
	Count(*sql.DB) (int, error)
	Save(*sql.DB) error
	GetByName(*sql.DB, string) error
}
