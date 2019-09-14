package consts

import (
	"errors"
)

const (
	SmtWhenWrongResp         = "Something when wrong"
	NameIsEmptyResp          = "Name is empty"
	BadInputDataResp         = "Bad input data"
	IdIsEmptyResp            = "Id is empty"
	IdIsNotNumberResp        = "Id is not number"
	InputDataInvalidResp     = "Invalid input data"
	NotSavedResp             = "Save error"
	NotExistResp             = "Not exist"
	PageAndCountRequiredResp = "Page and count required"
	PageMustNumber           = "Page must be a number"
	CountMustNumber          = "Count must be a number"
	NotUpdatedResp           = "Not updated"
	NotDeletedResp           = "Not deleted"
)

var (
	ErrNotFound = errors.New("not exist")
)
