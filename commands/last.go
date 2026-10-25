package command

import (
	"fmt"
	"regexp"
	"strings"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// LastCommand definition
type LastCommand struct {
	Db db.ZbotDatabase
}

// ShowHelp returns the help message for the command.
func (handler *LastCommand) ShowHelp(args ...string) string { return LastHelpMessage }
func (handler *LastCommand) ShowDescription() string        { return LastDescription }

const LastHelpMessage = "!last - get the last learned term"
const LastDescription = "Get the last learned term."
const LAST_ITEMS = 10

// SetDb set db connection if the module need it
func (handler *LastCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *LastCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!last$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		lastItems, err := handler.Db.Last(uint(channelID), LAST_ITEMS)
		if err != nil {
			slog.Error("last query error", "err", err)
			return "", err
		}
		return PrintTerms(lastItems), nil
	}
	return "", ErrNextCommand
}

// PrintTerms .
func PrintTerms(items []db.Definition) string {
	keys := make([]string, 0, len(items))
	for item := range items {
		keys = append(keys, items[item].Term)
	}
	return fmt.Sprintf("[ %s ]", strings.Join(keys, " "))
}
