package zbot

import (
	"log/slog"
	"slices"
)

// IsAdmin checks if the given username is an admin.
func (b *Bot) IsAdmin(username string) bool {
	slog.Debug("List of admins", "admins", b.cfg.Admins)
	slog.Debug("Checking admin status for user", "user", username)
	return slices.Contains(b.cfg.Admins, username)
}
