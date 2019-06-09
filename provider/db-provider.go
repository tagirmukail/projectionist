package provider

import (
	"database/sql"
	"fmt"
	"projectionist/models"
)

type DBProvider struct {
	db *sql.DB
}

func NewDBProvider(db *sql.DB) *DBProvider {
	return &DBProvider{
		db: db,
	}
}

func (p *DBProvider) Save(m interface{}) error {
	var model, ok = m.(models.Model)
	if !ok {
		return fmt.Errorf("this interface not model")
	}

	return model.Save(p.db)
}

func (p *DBProvider) GetByName(m interface{}, name string) error {
	model, ok := m.(models.Model)
	if !ok {
		return fmt.Errorf("this interface not model")
	}

	return model.GetByName(p.db, name)
}

func (p *DBProvider) IsExist(m interface{}) (error, bool) {
	model, ok := m.(models.Model)
	if !ok {
		return fmt.Errorf("this interface not model"), false
	}

	return model.IsExist(p.db)
}
