package provider

import (
	"database/sql"
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

func (p *DBProvider) Save(m models.Model) error {
	return m.Save(p.db)
}

func (p *DBProvider) GetByName(m models.Model, name string) error {
	return m.GetByName(p.db, name)
}

func (p *DBProvider) GetByID(m models.Model, id int64) error {
	return m.GetByID(p.db, id)
}

func (p *DBProvider) IsExist(m models.Model) (error, bool) {
	return m.IsExist(p.db)
}
