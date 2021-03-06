package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

type Status int

const (
	_ = iota
	Alive
	Dead
)

type Service struct {
	ID        int     `json:"id";db:"id"`
	Name      string  `json:"name";db:"name"`
	Link      string  `json:"link";db:"link"`
	Token     string  `json:"token";db:"token"`
	Frequency int     `json:"frequency";db:"frequency"`
	Status    Status  `json:"status";db:"status"`
	Deleted   int     `json:"deleted";db:"deleted"`
	Emails    []Email `json:"emails"`
}

func (s *Service) Validate() error {
	if s.Name == "" || len(s.Name) > 255 {
		return fmt.Errorf("invalid service name")
	}

	if s.Link == "" || len(s.Link) > 255 {
		return fmt.Errorf("invalid service link")
	}

	u, err := url.ParseRequestURI(s.Link)
	if err != nil {
		return fmt.Errorf("link %s error: %v", s.Link, err)
	}

	if u.Host == "" {
		return fmt.Errorf("link %s not valid", s.Link)
	}

	if s.Token == "" || len(s.Token) > 255 {
		return fmt.Errorf("invalid service token")
	}

	if s.Frequency == 0 {
		return fmt.Errorf("frequency must be not 0")
	}

	if s.Status != Alive && s.Status != Dead {
		return fmt.Errorf("invalid service status: %v", s.Status)
	}

	for _, email := range s.Emails {
		err = email.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) IsExistByName(db *sql.DB) (error, bool) {
	var serviceName string

	var err = db.QueryRow("SELECT name FROM services where name=?", s.Name).Scan(&serviceName)
	if err != nil && err != sql.ErrNoRows {
		return err, false
	}

	if serviceName == "" {
		return nil, false
	}

	return nil, true
}
func (s *Service) Count(db *sql.DB) (int, error) {
	var count int
	var err = db.QueryRow("SELECT count(id) FROM services").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (s *Service) Save(db *sql.DB) error {
	result, err := db.Exec(
		"INSERT INTO services (name, link, Token, frequency, status) VALUES (?,?,?,?,?)",
		s.Name,
		s.Link,
		s.Token,
		s.Frequency,
		s.Status,
	)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(lastInsertID)

	return nil
}

func (s *Service) GetByName(db *sql.DB, name string) error {
	return db.QueryRow(
		"SELECT id, name, link, Token, frequency, status, deleted FROM services WHERE name=?", name).Scan(
		&s.ID,
		&s.Name,
		&s.Link,
		&s.Token,
		&s.Frequency,
		&s.Status,
		&s.Deleted,
	)
}

func (s *Service) GetByID(db *sql.DB, id int64) error {
	return db.QueryRow(
		"SELECT id, name, link, Token, frequency, status, deleted FROM services WHERE id=?", id).Scan(
		&s.ID,
		&s.Name,
		&s.Link,
		&s.Token,
		&s.Frequency,
		&s.Status,
		&s.Deleted,
	)
}

func (s *Service) Pagination(db *sql.DB, start, end int) ([]Model, error) {
	var result []Model

	raws, err := db.Query(
		"SELECT id, name, link, Token, frequency, status, deleted FROM services ORDER BY id ASC limit ?, ?",
		start, end)
	if err != nil {
		return nil, err
	}

	for raws.Next() {
		var service = &Service{}
		err = raws.Scan(
			&service.ID,
			&service.Name,
			&service.Link,
			&service.Token,
			&service.Frequency,
			&service.Status,
			&service.Deleted,
		)
		if err != nil {
			return nil, err
		}

		service.Emails, err = service.GetEmails(db)
		if err != nil {
			return nil, err
		}

		result = append(result, service)
	}

	return result, nil
}

func (s *Service) Update(db *sql.DB, id int) error {
	var query, args = s.buildServiceUpdateQuery(id)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var res sql.Result
	if len(args) > 1 {
		res, err = db.Exec(
			query, args...,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("update query: %s error: %v", query, err)
		}

		rowsCount, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}

		if rowsCount == 0 {
			tx.Rollback()
			return fmt.Errorf("service with id %d not updated", id)
		}
	}

	for _, email := range s.Emails {
		err = email.Validate()
		if err != nil {
			tx.Rollback()
			return err
		}

		if email.ServiceID != id {
			tx.Rollback()
			return fmt.Errorf("invalid service id for email %v", email.Email)
		}
		err = email.Save(db)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Service) Delete(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	res, err := db.Exec("UPDATE services SET deleted=1 WHERE id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsCount == 0 {
		tx.Rollback()
		return fmt.Errorf("service with id %d not deleted", id)
	}

	res, err = db.Exec("DELETE FROM emails WHERE service_id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Service) GetID() int {
	return s.ID
}

func (s *Service) SetID(id int) {
	s.ID = id
}

func (s *Service) SetName(name string) {
	s.Name = name
}

func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) SetDeleted() {
	s.Deleted = 1
}

func (s *Service) IsDeleted() bool {
	return s.Deleted > 0
}

func (s *Service) GetEmails(db *sql.DB) ([]Email, error) {
	s.Emails = []Email{}

	rows, err := db.Query("SELECT id, service_id, email FROM emails WHERE service_id=?", s.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var email = Email{}
		err = rows.Scan(
			&email.ID,
			&email.ServiceID,
			&email.Email,
		)
		if err != nil {
			return nil, err
		}

		s.Emails = append(s.Emails, email)
	}

	return s.Emails, nil
}

func (s *Service) buildServiceUpdateQuery(id int) (string, []interface{}) {
	var queryBuild = strings.Builder{}

	var args []interface{}

	queryBuild.WriteString(`UPDATE services SET `)

	if s.Name != "" {
		queryBuild.WriteString(`name=?, `)
		args = append(args, s.Name)
	}

	if s.Link != "" {
		queryBuild.WriteString(`link=?, `)
		args = append(args, s.Link)
	}

	if s.Token != "" {
		queryBuild.WriteString(`token=?, `)
		args = append(args, s.Token)
	}

	if s.Frequency > 0 {
		queryBuild.WriteString(`frequency=?, `)
		args = append(args, s.Frequency)
	}

	if s.Status == Alive || s.Status == Dead {
		queryBuild.WriteString(`status=?, `)
		args = append(args, s.Status)
	}

	query := strings.TrimRight(queryBuild.String(), ", ")

	queryBuild.Reset()
	queryBuild.WriteString(query)

	queryBuild.WriteString(` WHERE id=?`)
	args = append(args, id)

	return queryBuild.String(), args
}
