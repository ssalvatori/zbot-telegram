package db

import "errors"

var (
	//ErrNotFound definition not found in db
	ErrNotFound = errors.New("Definition not found")
	//ErrLocked definition is locked
	ErrLocked = errors.New("Definition locked")
	//ErrInternalError internal error
	ErrInternalError = errors.New("Internal error")
)
