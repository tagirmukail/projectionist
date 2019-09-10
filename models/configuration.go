package models

import "fmt"

type Configuration struct {
	ID      int                    `json:"id"`
	Name    string                 `json:"username"`
	Config  map[string]interface{} `json:"config"`
	Deleted int                    `json:"deleted"`
}

// TODO: implement Save, Update, Delete methods

func (c *Configuration) SetDBCtx(interface{}) error {

	return nil
}

func (c *Configuration) Validate() error {
	if c.ID == 0 {
		return fmt.Errorf("id must be not 0")
	}
	return nil
}

func (c *Configuration) IsExistByName() (error, bool) {
	return nil, false
}

func (c *Configuration) Count() (int, error) {
	return 0, nil
}

func (c *Configuration) Save() error {

	return nil
}

func (c *Configuration) GetByName(name string) error {
	return nil
}

func (c *Configuration) GetByID(id int64) error {
	return nil
}

func (c *Configuration) Pagination(start, stop int) ([]Model, error) {
	return nil, nil
}

func (c *Configuration) Update(id int) error {
	return nil
}

func (c *Configuration) Delete(id int) error {
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

func (c *Configuration) SetDeleted() {
	c.Deleted = 1
}
