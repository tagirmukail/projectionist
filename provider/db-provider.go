package provider

import (
	"database/sql"
	"projectionist/models"
	"time"
)

const writeTimeout = 100 * time.Millisecond

type DBProvider struct {
	db *sql.DB
}

func NewDBProvider(db *sql.DB) *DBProvider {
	return &DBProvider{
		db: db,
	}
}

func (p *DBProvider) Save(m models.Model) error {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return err
	}

	return saveProcessErrBusy(m.Save)
}

func (p *DBProvider) GetByName(m models.Model, name string) (models.Model, error) {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return nil, err
	}

	return m, m.GetByName(name)
}

func (p *DBProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return nil, err
	}

	return m, m.GetByID(id)
}

func (p *DBProvider) IsExistByName(m models.Model) (error, bool) {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return err, false
	}

	return m.IsExistByName()
}

func (p *DBProvider) Count(m models.Model) (int, error) {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return 0, err
	}

	return m.Count()
}

func (p *DBProvider) Pagination(m models.Model, start, stop int) ([]models.Model, error) {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return nil, err
	}

	return m.Pagination(start, stop)
}

func (p *DBProvider) Update(m models.Model, id int) error {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return err
	}

	return processErrBusy(id, m.Update)
}

func (p *DBProvider) Delete(m models.Model, id int) error {
	err := m.SetDBCtx(p.db)
	if err != nil {
		return err
	}

	return processErrBusy(id, m.Delete)
}
