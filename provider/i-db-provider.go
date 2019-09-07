package provider

import "projectionist/models"

type IDBProvider interface {
	Save(models.Model) error
	GetByID(models.Model, int64) (models.Model, error)
	GetByName(models.Model, string) error
	IsExistByName(models.Model) (error, bool)
	Count(models.Model) (int, error)
	Pagination(models.Model, int, int) ([]models.Model, error)
	Update(models.Model, int) error
	Delete(models.Model, int) error
}
