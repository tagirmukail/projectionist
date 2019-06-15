package provider

import "projectionist/models"

type IDBProvider interface {
	Save(models.Model) error
	GetByID(models.Model, int64) error
	GetByName(models.Model, string) error
	IsExist(models.Model) (error, bool)
}
