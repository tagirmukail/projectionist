package errors

var (
	ErrAlreadyExist = New(
		1,
		500,
		"already exist",
		"model with this key already exist")
	ErrNotExist = New(
		2,
		404,
		"not exist",
		"model with this key not exist",
	)
	ErrUnknownModel = New(
		3,
		500,
		"something when wrong",
		"model with this type not unknown")
)
