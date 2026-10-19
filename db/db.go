package db

import (
	"database/sql"

	"gorm.io/gorm"
)

// Definition struct
type Definition struct {
	// gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Term      string `gorm:"uniqueIndex:idx_term_channel"`
	Meaning   string
	Author    string
	ChannelID uint `gorm:"uniqueIndex:idx_term_channel"`
	Hits      uint
	Link      sql.NullInt32

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

// Channel is a Telegram chat the bot operates in.
type Channel struct {
	ID              uint  `gorm:"primaryKey"`
	TelegramID      int64 `gorm:"uniqueIndex"`
	Name            string
	LearnEnabled    bool  `gorm:"default:true"`
	ExternalEnabled bool  `gorm:"default:true"`
	CreatedAt       int64 `gorm:"autoCreateTime"`
}

// ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	// Engine returns the backend name (e.g. "sqlite")
	Engine() string
	GetConnectionInfo() string
	Init() error
	Close()
	// Definition methods are scoped by Channel.ID.
	Statistics(channelID uint) (string, error)
	Last(channelID uint, limit int) ([]Definition, error)
	Append(Definition) error
	Top(channelID uint, limit int) ([]Definition, error)
	Rand(channelID uint, limit int) ([]Definition, error)
	Get(term string, channelID uint) (Definition, error)
	Set(Definition) (string, error)
	_set(string, Definition) error
	Find(criteria string, channelID uint, limit int) ([]Definition, error)
	Search(criteria string, channelID uint, limit int) ([]Definition, error)
	Forget(Definition) error
	UserLevel(string) (string, error)
	//UserIgnoreInsert ignore username for durationSeconds, recording who requested it and on which chat
	UserIgnoreInsert(string, int, string, uint) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string, uint) bool
	UserCleanupIgnorelist() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(Definition) error
	IncreaseHits(uint) error

	// Channel methods are keyed by the Telegram chat ID.
	AddChannel(telegramID int64, name string) error
	RemoveChannel(telegramID int64) (bool, error)
	ListChannels() ([]Channel, error)
	GetChannel(telegramID int64) (Channel, error)
	// DisableLearn(string) (string, error) // Disable learning for a channel
	// IsLearnEnabled(string) (bool, error) // Check if learning is enabled for a channel
}
