package models

const (
	NotDeleted int = iota
	Deleted
)

type Model interface {
	Validate() error
	GetID() int
	SetID(int)
	SetName(name string)
	GetName() string
	SetDeleted()
	IsDeleted() bool
}
