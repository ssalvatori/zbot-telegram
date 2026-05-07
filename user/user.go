// Package user provides user identity extraction and permission level management
// for the zbot Telegram bot. It bridges Telegram sender information with the bot's
// internal permission system, allowing commands to enforce access control based on
// numeric user levels stored in the database.
//
// A User is constructed from a Telegram message sender via [BuildUser], which
// normalizes the username and retrieves the current permission level. Commands
// then call [User.IsAllow] to check if the user meets the minimum level required
// to execute.
//
// Levels are integers where higher values grant more privileges:
//   - 0: default (all users)
//   - 100: moderator commands (e.g., ignore)
//   - 1000: admin commands (e.g., lock, forget)
package user

import (
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/ssalvatori/zbot-telegram/db"

	"strconv"

	"log/slog"
)

// User represents a bot user with identity and permission information.
// Username and Ident are always stored in lowercase for case-insensitive matching.
type User struct {
	Username string // Telegram username (lowercase), falls back to first name if unset
	Ident    string // Telegram first name (lowercase)
	Host     string // Reserved for future use
	Level    int    // Permission level from the database (0 = default)
}

// BuildUser creates a User from Telegram sender information, normalizing the
// username to lowercase and fetching the current permission level from the database.
// If the sender has no Telegram username set, the first name is used instead.
//
// Example:
//
//	user := BuildUser(message.From, database)
//	if !user.IsAllow(requiredLevel) {
//	    return "permission denied"
//	}
func BuildUser(sender *models.User, db db.ZbotDatabase) User {
	user := User{}
	user.Ident = strings.ToLower(sender.FirstName)

	if sender.Username != "" {
		user.Username = strings.ToLower(sender.Username)
	} else {
		user.Username = strings.ToLower(sender.FirstName)
	}

	user.Level = GetUserLevel(db, sender.Username)

	return user
}

// GetUserLevel retrieves the permission level for a user from the database.
// Returns 0 if the user has no level set or if a database error occurs.
func GetUserLevel(Db db.ZbotDatabase, username string) int {
	userLevel, err := Db.UserLevel(username)

	if err != nil {
		slog.Error("user level error", "err", err)
		return 0
	}

	userLevelInt, _ := strconv.Atoi(userLevel)

	return userLevelInt
}

// IsAllow reports whether the user's level is greater than or equal to the
// required level. Use this to gate command execution behind permission checks.
//
// Example:
//
//	if !user.IsAllow(1000) {
//	    return "you need admin privileges for this command"
//	}
func (u User) IsAllow(level int) bool {
	result := false

	if u.Level >= level {
		result = true
	}

	return result
}
