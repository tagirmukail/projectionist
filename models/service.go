package models

import (
	"database/sql"
	"fmt"
	"net/url"
)

type Status int

const (
	_ = iota
	Alive
	Dead
)

type Service struct {
	dbCtx     *sql.DB
	ID        int    `json:"id";db:"id"`
	Name      string `json:"name";db:"name"`
	Link      string `json:"link";db:"link"`
	Token     string `json:"token";db:"token"`
	Frequency int    `json:"frequency";db:"frequency"`
	Status    Status `json:"status";db:"status"`
	Deleted   int    `json:"deleted";db:"deleted"`
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

	return nil
}

func (s *Service) SetDBCtx(iDB interface{}) error {
	db, ok := iDB.(*sql.DB)
	if !ok {
		return fmt.Errorf("%v is not sql.DB", iDB)
	}

	s.dbCtx = db

	return nil
}
func (s *Service) IsExistByName() (error, bool) {
	var serviceName string

	var err = s.dbCtx.QueryRow("SELECT name FROM services where name=?", s.Name).Scan(&serviceName)
	if err != nil && err != sql.ErrNoRows {
		return err, false
	}

	if serviceName == "" {
		return nil, false
	}

	return nil, true
}
func (s *Service) Count() (int, error) {
	var count int
	var err = s.dbCtx.QueryRow("SELECT count(id) FROM services").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (s *Service) Save() error {
	result, err := s.dbCtx.Exec(
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

// TODO implement this method
func (s *Service) GetByName(name string) error { return nil }

func (s *Service) GetByID(id int64) error {
	return s.dbCtx.QueryRow(
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

func (s *Service) Pagination(start, end int) ([]Model, error) {
	var result []Model

	raws, err := s.dbCtx.Query(
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

		result = append(result, service)
	}

	return result, nil
}

// TODO implement this method
func (s *Service) Update(id int) error { return nil }

// TODO implement this method
func (s *Service) Delete(id int) error { return nil }

func (s *Service) GetID() int {
	return s.ID
}

// TODO implement this method
func (s *Service) SetID(id int) { return }

// TODO implement this method
func (s *Service) GetName() string { return "" }

// TODO implement this method
func (s *Service) SetDeleted() {}

// TODO implement this method
func (s *Service) IsDeleted() bool { return false }
