package db

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//Definition struct
type Definition struct {
	gorm.Model
	ID      uint
	Term    string
	Meaning string
	Author  string
	Date    string
	Chat    string
	Hits    uint
	Link    sql.NullInt32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime   `gorm:"default:null"`
	Locked    sql.NullBool   `gorm:"default:false"`
	LockedBy  sql.NullString `gorm:"default:null"`
	DeletedBy sql.NullString `gorm:"default:null"`
}

//UserIgnore .
type UserIgnore struct {
	Username string
	Since    string
	Until    string
}

//ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	GetConnectionInfo() string
	Init() error
	Close()
	Statistics() (string, error)
	Last(int) ([]Definition, error)
	Append(Definition) error
	Top(int) ([]Definition, error)
	Rand(int) ([]Definition, error)
	Get(string, string) (Definition, error)
	Set(Definition) (string, error)
	_set(string, Definition) error
	Find(string, string, int) ([]Definition, error)
	Search(string, string, int) ([]Definition, error)
	Forget(Definition) error
	UserLevel(string) (string, error)
	UserIgnoreInsert(string) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string) bool
	UserCleanupIgnorelist() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(Definition) error
	IncreaseHits(uint) error
}
