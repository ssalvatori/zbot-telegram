package db

import (
	"database/sql"

	"gorm.io/gorm"
)

// Definition struct
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

// UserIgnore .
type UserIgnore struct {
	Username   string
	CreatedAt  int64 `gorm:"autoCreateTime"`
	ValidUntil int64
	CreatedBy  string
	ChannelID  uint
}

// Channel
type Channel struct {
	ID              uint `gorm:"primaryKey"`
	Name            string
	LearnEnabled    bool  `grom:"default:true"`
	ExternalEnabled bool  `grom:"default:true"`
	CreatedAt       int64 `gorm:"autoCreateTime"`
}

// ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	// Engine returns the backend name (e.g. "sqlite")
	Engine() string
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
	//UserIgnoreInsert ignore username for durationSeconds, recording who requested it and on which chat
	UserIgnoreInsert(string, int, string, uint) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string, uint) bool
	UserCleanupIgnorelist() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(Definition, string) error
	IncreaseHits(uint) error

	// Channel
	AddChannel(string) error            // Add a new channel to the database
	RemoveChannel(string) (bool, error) // Remove channel from database
	ListChannels() ([]Channel, error)   // List all channels in database
	GetChannel(string) (Channel, error) // Get channel information
	// DisableLearn(string) (string, error) // Disable learning for a channel
	// IsLearnEnabled(string) (bool, error) // Check if learning is enabled for a channel
}
