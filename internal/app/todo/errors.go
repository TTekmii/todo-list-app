package todo

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNoFieldsToUpdate = errors.New("no field to update")
	ErrNotFound         = errors.New("not found")
)
