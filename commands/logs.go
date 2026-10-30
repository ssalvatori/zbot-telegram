package command

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/user"
)

const LogsHelpMessage = "!logs - show the last log entries (admins, private chat)"
const LogsDescription = "Reads the plain text log file and returns its last 10 lines."
const logsDefaultLimit = 10

// DefaultLogsPath is the fixed file the bot writes its logs to (see setupLog
// in main.go) and that the !logs command reads back.
const DefaultLogsPath = "zbot.log"

const telegramMaxMessageLength = 4000 // Telegram's maximum message length is 4096 characters, but we leave some room for the truncation message and other formatting.

type LogsCommand struct {
	Path string
}

func (handler *LogsCommand) ShowHelp(args ...string) string { return LogsHelpMessage }
func (handler *LogsCommand) ShowDescription() string        { return LogsDescription }

func (handler *LogsCommand) ProcessText(text string, u user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!logs$`)

	if isCommand || commandPattern.MatchString(text) {
		msg, err := handler.readLogs()
		if err != nil {
			slog.Error("read logs error", "err", err)
		}
		slog.Info(fmt.Sprintf("Fetching %d log entries requested by %s", logsDefaultLimit, u.Username))
		return truncateMessage(msg, telegramMaxMessageLength), nil
	}

	return "", ErrNextCommand
}

// truncateMessage truncates a message to the specified maximum length, appending an ellipsis if truncation occurs.
func truncateMessage(msg string, maxLength int) string {
	if len(msg) > maxLength {
		msg = msg[len(msg)-maxLength:]
		msg = msg + "... (truncated)"
	}
	return msg
}

// readLogs reads the last log entries from the log file and returns them as an string
func (handler *LogsCommand) readLogs() (string, error) {
	// TODO: mask sensitive values (e.g. the Telegram token or database
	// connection strings) before returning the output, so !logs cannot
	// leak them in plaintext.
	data, err := os.ReadFile(handler.Path)
	if err != nil {
		return "No log entries yet.", nil
	}

	lines := strings.Split(string(data), "\n")
	// Drop the trailing empty element produced by a file that ends in "\n".
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	entries := make([]string, 0, logsDefaultLimit)
	for _, l := range lines {
		if l == "" {
			continue
		}
		entries = append(entries, l)
	}
	if len(entries) == 0 {
		return "No log entries yet.", nil
	}
	if len(entries) > logsDefaultLimit {
		entries = entries[len(entries)-logsDefaultLimit:]
	}

	return strings.Join(entries, "\n"), nil
}
