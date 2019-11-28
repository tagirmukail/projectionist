package models

import (
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

	if _, err := url.ParseRequestURI(s.Link); err != nil {
		return fmt.Errorf("link %s error: %v", s.Link, err)
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
