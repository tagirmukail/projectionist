package models

const (
	NotDeleted int = iota
	Deleted
)

type Model interface {
	SetDBCtx(interface{}) error
	Validate() error
	IsExistByName() (error, bool)
	Count() (int, error)
	Save() error
	GetByName(string) error
	GetByID(int64) error
	Pagination(int, int) ([]Model, error)
	Update(int) error
	Delete(int) error
	GetID() int
	SetID(int)
	GetName() string
	SetDeleted()
}
