package db

import (
	"strings"
	"time"

	"log/slog"
)

// UserIgnoreList return users currently on the ignore list
func (d *ZbotDatabaseSqlite) UserIgnoreList() ([]UserIgnore, error) {
	ignoredUsers := []UserIgnore{}
	now := time.Now().Unix()
	if err := d.DB.Model(&UserIgnore{}).Where("valid_until > ?", now).Order("created_at desc").Find(&ignoredUsers).Error; err != nil {
		slog.Error("ignore list query error", "err", err)
		return nil, err
	}
	return ignoredUsers, nil
}

// UserLevel Mock
func (d *ZbotDatabaseSqlite) UserLevel(str string) (string, error) {
	return "bnil", nil
}

// UserCheckIgnore return true if the user is currently on the ignore list
func (d *ZbotDatabaseSqlite) UserCheckIgnore(username string, channelID uint) bool {
	var count int64
	now := time.Now().Unix()
	if err := d.DB.Model(&UserIgnore{}).Where("username = ? COLLATE NOCASE", username).Where("channel_id = ?", channelID).Where("valid_until > ?", now).Count(&count).Error; err != nil {
		slog.Error("check ignore query error", "err", err)
		return false
	}
	return count > 0
}

// UserIgnoreInsert Add a user to ignore list
func (d *ZbotDatabaseSqlite) UserIgnoreInsert(username string, durationSeconds int, createdBy string, channelID uint) error {
	ignore := UserIgnore{
		Username:   strings.ToLower(username),
		ValidUntil: time.Now().Add(time.Duration(durationSeconds) * time.Second).Unix(),
		CreatedBy:  createdBy,
		ChannelID:  channelID,
	}

	if err := d.DB.Create(&ignore).Error; err != nil {
		slog.Error("ignore insert error", "err", err)
		return err
	}
	return nil
}

// UserCleanupIgnorelist Cleanup ignore list
func (d *ZbotDatabaseSqlite) UserCleanupIgnorelist() error {
	now := time.Now().Unix()
	if err := d.DB.Where("valid_until <= ?", now).Delete(&UserIgnore{}).Error; err != nil {
		slog.Error("ignore list cleanup error", "err", err)
		return err
	}
	return nil
}
