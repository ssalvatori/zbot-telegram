package db

import "testing"

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
