package provider

import (
	"github.com/mattn/go-sqlite3"
	"time"
)

func saveProcessErrBusy(save func() error) error {
	err := save()
	for i := 0; i < 5; i++ {
		if err != sqlite3.ErrBusy {
			return err
		}
		time.Sleep(writeTimeout)

		err = save()
	}

	return err
}

func processErrBusy(id int, f func(id int) error) error {
	err := f(id)
	for i := 0; i < 5; i++ {
		if err != sqlite3.ErrBusy {
			return err
		}
		time.Sleep(writeTimeout)

		err = f(id)
	}

	return err
}
