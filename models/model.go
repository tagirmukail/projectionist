package models

import "database/sql"

const (
	NotDeleted int = iota
	Deleted
)

type Model interface {
	Validate() error
	IsExistByName(*sql.DB) (error, bool)
	Count(*sql.DB) (int, error)
	Save(*sql.DB) error
	GetByName(*sql.DB, string) error
	GetByID(*sql.DB, int64) error
	Pagination(*sql.DB, int, int) ([]Model, error)
	Update(*sql.DB, int) error
	Delete(*sql.DB, int) error
}
