package models

import (
	"database/sql"
	"fmt"
	"regexp"
)

var RegexEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Email struct {
	ID        int    `json:"id"`
	ServiceID int    `json:"service_id"`
	Email     string `json:"email"`
}

func (e *Email) Validate() error {
	if e.Email == "" {
		return fmt.Errorf("empty email")
	}

	if e.ServiceID == 0 {
		return fmt.Errorf("for email %v service id is empty", e.Email)
	}

	if !RegexEmail.MatchString(e.Email) {
		return fmt.Errorf("is not email: %s", e.Email)
	}

	return nil
}

func (e *Email) IsExistByName(db *sql.DB) (error, bool) { return nil, false }
func (e *Email) Count(db *sql.DB) (int, error)          { return 0, nil }

func (e *Email) Save(db *sql.DB) error {
	return insertEmail(db, *e)
}

func (e *Email) GetByName(db *sql.DB, name string) error                    { return nil }
func (e *Email) GetByID(db *sql.DB, id int64) error                         { return nil }
func (e *Email) Pagination(db *sql.DB, start int, end int) ([]Model, error) { return nil, nil }
func (e *Email) Update(db *sql.DB, id int) error                            { return nil }
func (e *Email) Delete(db *sql.DB, id int) error {
	res, err := db.Exec("DELETE FROM emails WHERE id=?", id)
	if err != nil {
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return fmt.Errorf("email with id %d not deleted", id)
	}

	return nil
}
func (e *Email) GetID() int          { return e.ID }
func (e *Email) SetID(id int)        { e.ID = id }
func (e *Email) GetName() string     { return e.Email }
func (e *Email) SetName(name string) { e.Email = name }
func (e *Email) SetDeleted()         {}
func (e *Email) IsDeleted() bool     { return false }

func insertEmail(db *sql.DB, email Email) error {
	var existEmail Email
	err := db.QueryRow("SELECT id, service_id, email FROM emails WHERE service_id=? AND email=?", email.ServiceID, email.Email).Scan(
		&existEmail.ID,
		&existEmail.ServiceID,
		&existEmail.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existEmail.ID != 0 {
		return fmt.Errorf("this email %s for service %d exist", existEmail.Email, existEmail.ServiceID)
	}

	res, err := db.Exec(`INSERT INTO emails (service_id, email) VALUES (?,?)`, email.ServiceID, email.Email)
	if err != nil {
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if int(rowsCount) == 0 {
		return fmt.Errorf("email %s not inserted", email.Email)
	}

	return nil
}
