package db

import (
	"database/sql"
)

// createTableUsers create table users if not exist
func createTableUsers(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(CREATE_TBL_USERS)
	return err
}

// createTableServices create table services if not exist
func createTableServices(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(CREATE_TBL_SERVICE)
	return err
}

func createTableEmails(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(CREATE_TBL_EMAILS)
	return err
}

func InitTables(sqlDB *sql.DB) error {
	var err error
	if err = createTableServices(sqlDB); err != nil {
		return err
	}

	if err = createTableUsers(sqlDB); err != nil {
		return err
	}

	if err = createTableEmails(sqlDB); err != nil {
		return err
	}

	return nil
}
