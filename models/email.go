package models

import (
	"database/sql"
	"fmt"
	"regexp"
)

var RegexEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Email struct {
	dbCtx     *sql.DB
	ID        int    `json:"id"`
	ServiceID int    `json:"service_id"`
	Email     string `json:"email"`
}

func (e *Email) SetDBCtx(iDB interface{}) error {
	db, ok := iDB.(*sql.DB)
	if !ok {
		return fmt.Errorf("%v is not sql.DB", iDB)
	}

	e.dbCtx = db

	return nil
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

func (e *Email) IsExistByName() (error, bool) { return nil, false }
func (e *Email) Count() (int, error)          { return 0, nil }

func (e *Email) Save() error {
	return insertEmail(e.dbCtx, *e)
}

func (e *Email) GetByName(name string) error                    { return nil }
func (e *Email) GetByID(id int64) error                         { return nil }
func (e *Email) Pagination(start int, end int) ([]Model, error) { return nil, nil }
func (e *Email) Update(id int) error                            { return nil }
func (e *Email) Delete(id int) error {
	res, err := e.dbCtx.Exec("DELETE FROM emails WHERE id=?", id)
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
func (e *Email) GetID() int      { return e.ID }
func (e *Email) SetID(id int)    { e.ID = id }
func (e *Email) GetName() string { return e.Email }
func (e *Email) SetDeleted()     {}
func (e *Email) IsDeleted() bool { return false }

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
