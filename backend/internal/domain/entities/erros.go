package entities

import "errors"

var (
	// Entity Errors
	ErrNotFoundEntity      = errors.New("Entity Not Found")
	ErrInvalidEntity       = errors.New("Invalid Entity")
	ErrAlreadyExistsEntity = errors.New("Entity Already Exists")
)
