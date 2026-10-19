package db

import "errors"

var (
	//ErrNotFound definition not found in db
	ErrNotFound = errors.New("Definition not found")
	//ErrLocked definition is locked
	ErrLocked = errors.New("Definition locked")
	//ErrInternalError internal error
	ErrInternalError = errors.New("internal error")
	//ErrLearnDisabledChannel learn commands are disabled for a channel
	ErrLearnDisabledChannel = errors.New("learn commands disabled for this channel")
	//ErrNoEngine no database engine selected in the configuration
	ErrNoEngine = errors.New("no database engine selected")
	//ErrUnsupportedEngine the configured database engine is not implemented
	ErrUnsupportedEngine = errors.New("unsupported database engine")
	//ErrTermCollisionLimit no unique term could be found within the retry limit
	ErrTermCollisionLimit = errors.New("term collision retry limit reached")
)
