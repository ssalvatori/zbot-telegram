package db

import "errors"

var (
	//ErrNotFound definition not found in db
	ErrNotFound = errors.New("Definition not found")
	//ErrLocked definition is locked
	ErrLocked = errors.New("Definition locked")
	//ErrInternalError internal error
	ErrInternalError = errors.New("Internal error")
	//ErrLearnDisabledChannel learn commands are disabled for a channel
	ErrLearnDisabledChannel = errors.New("Learn commands disabled for this channel")
)
