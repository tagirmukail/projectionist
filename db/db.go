package db

import (
	"database/sql"
)

// CreateTableUsers create table users if not exist
func CreateTableUsers(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(CREATE_TBL_USERS)
	return err
}

// CreateTableServices create table services if not exist
func CreateTableServices(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(CREATE_TBL_SERVICE)
	return err
}
