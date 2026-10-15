package zbot

import "slices"

// IsAdmin checks if the given username is an admin.
func (b *Bot) IsAdmin(username string) bool {
	return slices.Contains(b.cfg.Admins, username)
}
