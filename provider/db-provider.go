package provider

import (
	"database/sql"
	"fmt"
	"projectionist/models"
	"projectionist/utils/errors"
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

func (p *DBProvider) GetDB() interface{} {
	return p.db
}

func (p *DBProvider) Save(m models.Model) error {
	switch m.(type) {
	case *models.Service:
		s := m.(*models.Service)
		return saveProcessErrBusy(p.db, s.Save)
	case *models.User:
		u := m.(*models.User)
		return saveProcessErrBusy(p.db, u.Save)
	case *models.Email:
		e := m.(*models.Email)
		return saveProcessErrBusy(p.db, e.Save)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return err
	}
}

func (p *DBProvider) GetByName(m models.Model, name string) (models.Model, error) {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return u, u.GetByName(p.db, name)
	case *models.Service:
		s := m.(*models.Service)
		return s, s.GetByName(p.db, name)
	case *models.Email:
		e := m.(*models.Email)
		return e, e.GetByName(p.db, name)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return nil, err
	}
}

func (p *DBProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return u, u.GetByID(p.db, id)
	case *models.Service:
		s := m.(*models.Service)
		return s, s.GetByID(p.db, id)
	case *models.Email:
		e := m.(*models.Email)
		return e, e.GetByID(p.db, id)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return nil, err
	}
}

func (p *DBProvider) IsExistByName(m models.Model) (error, bool) {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return u.IsExistByName(p.db)
	case *models.Service:
		s := m.(*models.Service)
		return s.IsExistByName(p.db)
	case *models.Email:
		e := m.(*models.Email)
		return e.IsExistByName(p.db)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return err, false
	}
}

func (p *DBProvider) Count(m models.Model) (int, error) {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return u.Count(p.db)
	case *models.Service:
		s := m.(*models.Service)
		return s.Count(p.db)
	case *models.Email:
		e := m.(*models.Email)
		return e.Count(p.db)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return 0, err
	}
}

func (p *DBProvider) Pagination(m models.Model, start, stop int) ([]models.Model, error) {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return u.Pagination(p.db, start, stop)
	case *models.Service:
		s := m.(*models.Service)
		return s.Pagination(p.db, start, stop)
	case *models.Email:
		e := m.(*models.Email)
		return e.Pagination(p.db, start, stop)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return nil, err
	}
}

func (p *DBProvider) Update(m models.Model, id int) error {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return processErrBusy(p.db, id, u.Update)
	case *models.Service:
		s := m.(*models.Service)
		return processErrBusy(p.db, id, s.Update)
	case *models.Email:
		e := m.(*models.Email)
		return processErrBusy(p.db, id, e.Update)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return err
	}
}

func (p *DBProvider) Delete(m models.Model, id int) error {
	switch m.(type) {
	case *models.User:
		u := m.(*models.User)
		return processErrBusy(p.db, id, u.Delete)
	case *models.Service:
		s := m.(*models.Service)
		return processErrBusy(p.db, id, s.Delete)
	case *models.Email:
		e := m.(*models.Email)
		return processErrBusy(p.db, id, e.Delete)
	default:
		err := errors.ErrUnknownModel
		err.SetArgs("type", fmt.Sprintf("%T", m))
		return err
	}
}
