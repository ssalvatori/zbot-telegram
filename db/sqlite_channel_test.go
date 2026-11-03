package db

import (
	"testing"
	"time"
)

func TestSqliteChannelOperations(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	const telegramID int64 = -100

	if err := d.AddChannel(telegramID, "testchat"); err != nil {
		t.Fatalf("AddChannel(%d, %q) error = %v", telegramID, "testchat", err)
	}

	channel, err := d.GetChannel(telegramID)
	if err != nil {
		t.Fatalf("GetChannel(%d) error = %v", telegramID, err)
	}
	if channel.ID == 0 || channel.TelegramID != telegramID || channel.Name != "testchat" || !channel.LearnEnabled || !channel.ExternalEnabled {
		t.Errorf("GetChannel(%d) = %+v, want stored channel with both features enabled", telegramID, channel)
	}

	if _, err := d.GetChannel(-999); err == nil {
		t.Errorf("GetChannel(-999) error = nil, want error for unknown channel")
	}

	channels, err := d.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels() error = %v", err)
	}
	if len(channels) != 1 || channels[0].TelegramID != telegramID {
		t.Errorf("ListChannels() = %+v, want one channel with TelegramID %d", channels, telegramID)
	}

	removed, err := d.RemoveChannel(telegramID)
	if err != nil {
		t.Fatalf("RemoveChannel(%d) error = %v", telegramID, err)
	}
	if !removed {
		t.Errorf("RemoveChannel(%d) = false, want true", telegramID)
	}

	channels, err = d.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels() after removal error = %v", err)
	}
	if len(channels) != 0 {
		t.Errorf("ListChannels() after removal = %+v, want no channels", channels)
	}
}

func TestSqliteDisableChannel(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	const telegramID int64 = -200

	if err := d.AddChannel(telegramID, "disableme"); err != nil {
		t.Fatalf("AddChannel(%d) error = %v", telegramID, err)
	}

	before := time.Now().Unix()
	disabled, err := d.DisableChannel(telegramID)
	if err != nil {
		t.Fatalf("DisableChannel(%d) error = %v", telegramID, err)
	}
	if !disabled {
		t.Fatalf("DisableChannel(%d) = false, want true", telegramID)
	}
	after := time.Now().Unix()

	ch, err := d.GetChannel(telegramID)
	if err != nil {
		t.Fatalf("GetChannel(%d) after disable error = %v", telegramID, err)
	}
	if !ch.Deactivate {
		t.Errorf("Deactivate = false, want true after DisableChannel")
	}
	if ch.DeactivatedAt < before || ch.DeactivatedAt > after {
		t.Errorf("DeactivatedAt = %d, want between %d and %d", ch.DeactivatedAt, before, after)
	}
}

func TestSqliteEnableChannel(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	const telegramID int64 = -300

	if err := d.AddChannel(telegramID, "enableme"); err != nil {
		t.Fatalf("AddChannel(%d) error = %v", telegramID, err)
	}
	if _, err := d.DisableChannel(telegramID); err != nil {
		t.Fatalf("DisableChannel(%d) error = %v", telegramID, err)
	}

	enabled, err := d.EnableChannel(telegramID)
	if err != nil {
		t.Fatalf("EnableChannel(%d) error = %v", telegramID, err)
	}
	if !enabled {
		t.Fatalf("EnableChannel(%d) = false, want true", telegramID)
	}

	ch, err := d.GetChannel(telegramID)
	if err != nil {
		t.Fatalf("GetChannel(%d) after enable error = %v", telegramID, err)
	}
	if ch.Deactivate {
		t.Errorf("Deactivate = true, want false after EnableChannel")
	}
	if ch.DeactivatedAt != 0 {
		t.Errorf("DeactivatedAt = %d, want 0 after EnableChannel", ch.DeactivatedAt)
	}
}

func TestSqliteCheckIfChannelExist(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	if d.CheckIfChannelExist(-404) {
		t.Errorf("CheckIfChannelExist(-404) = true, want false for unknown channel")
	}

	const telegramID int64 = -400
	if err := d.AddChannel(telegramID, "exists"); err != nil {
		t.Fatalf("AddChannel(%d) error = %v", telegramID, err)
	}

	if !d.CheckIfChannelExist(telegramID) {
		t.Errorf("CheckIfChannelExist(%d) = false, want true after AddChannel", telegramID)
	}
}
