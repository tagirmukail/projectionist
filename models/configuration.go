package models

import (
	"fmt"
)

const (
	FileIDNAMEPtrn = "%d|%s.json"
	FilePath       = "%s/%s.json"
	FilePathPtrn   = "%s/" + FileIDNAMEPtrn
)

type Configuration struct {
	ID      int                    `json:"id"`
	Name    string                 `json:"name"`
	Config  map[string]interface{} `json:"config"`
	Deleted int                    `json:"deleted"`
}

func (c *Configuration) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name must be not empty")
	}

	if c.Config == nil || len(c.Config) == 0 {
		return fmt.Errorf("config field must be not empty")
	}

	return nil
}

func (c *Configuration) GetID() int {
	return c.ID
}

func (c *Configuration) SetID(id int) {
	c.ID = id
}

func (c *Configuration) GetName() string {
	return c.Name
}

func (c *Configuration) SetName(name string) {
	c.Name = name
}

func (c *Configuration) SetDeleted() {
	c.Deleted = 1
}

func (c *Configuration) IsDeleted() bool {
	return c.Deleted > 0
}
