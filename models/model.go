package models

import "database/sql"

type Model interface {
	Validate() error
	IsExist(*sql.DB) (error, bool)
	Save(*sql.DB) error
}
