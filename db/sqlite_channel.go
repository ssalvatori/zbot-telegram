package db

import (
	"log/slog"
	"time"
)

func (d *ZbotDatabaseSqlite) ListChannels() ([]Channel, error) {
	var channels []Channel
	if err := d.DB.Model(&Channel{}).Find(&channels).Error; err != nil {
		slog.Error("list channels error", "err", err)
		return nil, err
	}
	return channels, nil
}

// AddChannel add a new channel to the database
func (d *ZbotDatabaseSqlite) AddChannel(telegramID int64, name string) error {
	if err := d.DB.Model(&Channel{}).Create(&Channel{TelegramID: telegramID, Name: name, LearnEnabled: true, ExternalEnabled: true}).Error; err != nil {
		slog.Error("add channel error", "err", err)
		return err
	}
	return nil
}

// RemoveChannel delete a channel from the database
func (d *ZbotDatabaseSqlite) RemoveChannel(telegramID int64) (bool, error) {
	if err := d.DB.Where("telegram_id = ?", telegramID).Delete(&Channel{}).Error; err != nil {
		slog.Error("delete channel error", "err", err)
		return false, err
	}
	return true, nil
}

// GetChannel get a channel information from the database
func (d *ZbotDatabaseSqlite) GetChannel(telegramID int64) (Channel, error) {
	var ch Channel
	if err := d.DB.Model(&Channel{}).Where("telegram_id = ?", telegramID).First(&ch).Error; err != nil {
		slog.Error("get channel error", "err", err)
		return Channel{}, err
	}
	return ch, nil
}

// DisableChannel disables a channel in the database and records when it was deactivated.
func (d *ZbotDatabaseSqlite) DisableChannel(telegramID int64) (bool, error) {
	updates := map[string]any{
		"deactivate":     true,
		"deactivated_at": time.Now().Unix(),
	}
	if err := d.DB.Model(&Channel{}).Where("telegram_id = ?", telegramID).Updates(updates).Error; err != nil {
		slog.Error("disable channel error", "err", err)
		return false, err
	}
	return true, nil
}

// EnableChannel disables a channel in the database and records when it was deactivated.
func (d *ZbotDatabaseSqlite) EnableChannel(channelID int64) (bool, error) {
	updates := map[string]any{
		"deactivate":     false,
		"deactivated_at": 0,
	}
	if err := d.DB.Model(&Channel{}).Where("telegram_id = ?", channelID).Updates(updates).Error; err != nil {
		slog.Error("disable channel error", "err", err)
		return false, err
	}
	return true, nil
}

func (d *ZbotDatabaseSqlite) CheckIfChannelExist(telegramID int64) bool {
	var count int64
	if err := d.DB.Model(&Channel{}).Where("telegram_id = ?", telegramID).Count(&count).Error; err != nil {
		slog.Error("check if channel exist error", "err", err)
		return false
	}
	return count > 0
}
