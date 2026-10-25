package command

import (
	"regexp"

	"github.com/ssalvatori/zbot-telegram/user"
)

const PingHelpMessage = "!ping - check if the bot is running"
const PingDescription = "Check if the bot is running."

// PingCommand command definition
type PingCommand struct{}

// ShowHelp returns the help message for the command.
func (handler *PingCommand) ShowHelp(args ...string) string { return PingHelpMessage }
func (handler *PingCommand) ShowDescription() string        { return PingDescription }

// ProcessText run command
func (handler *PingCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!ping$`)

	if isCommand || commandPattern.MatchString(text) {
		return "pong!!", nil
	}

	return "", ErrNextCommand
}
