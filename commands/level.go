package command

import (
	"fmt"
	"regexp"
	"strconv"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"

	"strings"

	"github.com/ssalvatori/zbot-telegram/user"
)

const LevelHelpMessage = "!level <subcommand> <username> <level> - manage user levels in the channel"
const LevelDescription = "Manage user levels in the channel."

// LevelCommand definition
type LevelCommand struct {
	Db db.ZbotDatabase
}

// AddUser add level for a given user
func (handler *LevelCommand) AddUser(targetUser string, level int, responsibleUser string) (string, error) {
	return "not ready", nil
}

// DelUser delete level for a given user
func (handler *LevelCommand) DelUser(userToCheck string, user string) (string, error) {
	return "not ready", nil
}

// GetLevel return level for a given user
func (handler *LevelCommand) GetLevel(user string) (string, error) {
	level, err := handler.Db.UserLevel(user)
	if err != nil {
		slog.Error("get level error", "err", err)
		return "", err
	}
	return fmt.Sprintf("%s level %s", user, level), nil
}

// ParseCommand identifies the user and subcommand in a level command.
func (handler *LevelCommand) ParseCommand(cmd string, user string) map[string]string {

	res := make(map[string]string)
	parsed := strings.Fields(cmd)

	res["user"] = user
	res["subcommand"] = "get"
	res["level"] = "0"

	if len(parsed) == 4 {
		switch parsed[1] {
		case "add":
			res["subcommand"] = "add"
			res["user"] = parsed[2]
			res["level"] = parsed[3]
		default:

		}
	}

	if len(parsed) == 3 {
		switch parsed[1] {
		case "del":
			res["subcommand"] = "del"
			res["user"] = parsed[2]
			res["level"] = "0"
		default:

		}
	}

	return res
}

// ShowHelp returns the help message for the command.
func (handler *LevelCommand) ShowHelp(args ...string) string { return LevelHelpMessage }
func (handler *LevelCommand) ShowDescription() string        { return LevelDescription }

// ProcessText run command
func (handler *LevelCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!level(\s|$)(\S*)\s?(\S+)?\s?(\d+)?`)
	var result string
	var err error

	if commandPattern.MatchString(text) {
		parsedCmd := handler.ParseCommand(text, user.Username)
		switch parsedCmd["subcommand"] {
		case "add":
			level, _ := strconv.Atoi(parsedCmd["level"])
			result, err = handler.AddUser(parsedCmd["user"], level, user.Username)
		case "del":
			result, err = handler.DelUser(parsedCmd["user"], user.Username)
		default:
			result, err = handler.GetLevel(parsedCmd["user"])
		}
		if err != nil {
			return "", err
		}

		return result, nil
	}
	return "", ErrNextCommand
}
