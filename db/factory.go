package db

import (
	"fmt"

	"github.com/ssalvatori/zbot-telegram/configuration"
)

// Engine names supported by New.
const (
	EngineSqlite = "sqlite"
)

// New builds the database backend selected by conf. The returned database is
// not connected yet; call Init on it before use.
func New(conf configuration.ConfigurationDb) (ZbotDatabase, error) {
	switch conf.Engine {
	case "":
		return nil, ErrNoEngine
	case EngineSqlite:
		return &ZbotDatabaseSqlite{File: conf.File}, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedEngine, conf.Engine)
	}
}
