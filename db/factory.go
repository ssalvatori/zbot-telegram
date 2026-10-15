package db

import (
	"fmt"
)

// Config contains the settings needed to create a database backend.
type Config struct {
	Engine   string `yaml:"engine"`
	Name     string `yaml:"name"`
	File     string `yaml:"file"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Engine names supported by New.
const (
	EngineSqlite = "sqlite"
)

// New builds the database backend selected by conf. The returned database is
// not connected yet; call Init on it before use.
func New(conf Config) (ZbotDatabase, error) {
	switch conf.Engine {
	case "":
		return nil, ErrNoEngine
	case EngineSqlite:
		return &ZbotDatabaseSqlite{File: conf.File}, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedEngine, conf.Engine)
	}
}
