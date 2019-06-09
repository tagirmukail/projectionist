package provider

type IDBProvider interface {
	Save(interface{}) error
	GetByName(interface{}, string) error
	IsExist(interface{}) (error, bool)
}
