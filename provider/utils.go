package provider

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"time"
)

func saveProcessErrBusy(db *sql.DB, save func(db *sql.DB) error) error {
	err := save(db)
	for i := 0; i < 5; i++ {
		if err != sqlite3.ErrBusyRecovery && err != sqlite3.ErrBusySnapshot && err != sqlite3.ErrBusy {
			return err
		}
		time.Sleep(writeTimeout)

		err = save(db)
	}

	return err
}

func processErrBusy(db *sql.DB, id int, f func(db *sql.DB, id int) error) error {
	err := f(db, id)
	for i := 0; i < 5; i++ {
		if err != sqlite3.ErrBusyRecovery && err != sqlite3.ErrBusySnapshot && err != sqlite3.ErrBusy {
			return err
		}
		time.Sleep(writeTimeout)

		err = f(db, id)
	}

	return err
}
