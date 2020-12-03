package db

import (
	"database/sql"

	"gorm.io/gorm"
)

//Definition struct
type Definition struct {
	// gorm.Model
	ID      uint
	Term    string `gorm:"uniqueIndex"`
	Meaning string
	Author  string
	// Date    string
	Chat string
	Hits uint
	Link sql.NullInt32

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
	Locked    sql.NullBool `gorm:"default:false"`
	LockedBy  sql.NullString
	DeletedBy sql.NullString
}

//UserIgnore .
type UserIgnore struct {
	Username   string
	CreatedAt  int64 `gorm:"autoCreateTime"`
	ValidUntil int64
	CreatedBy  string
	Chat       string
}

//ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	GetConnectionInfo() string
	Init() error
	Close()
	Statistics(string) (string, error)
	Last(string, int) ([]Definition, error)
	Append(Definition, string) error
	Top(string, int) ([]Definition, error)
	Rand(string, int) ([]Definition, error)
	Get(string, string) (Definition, error)
	Set(Definition) (string, error)
	_set(string, Definition) error
	Find(string, string, int) ([]Definition, error)
	Search(string, string, int) ([]Definition, error)
	Forget(Definition, string) error
	UserLevel(string) (string, error)
	UserIgnoreInsert(string) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string) bool
	UserCleanupIgnorelist() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(Definition, string) error
	IncreaseHits(uint) error
}
