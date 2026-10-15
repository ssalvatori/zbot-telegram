package command

import (
	"regexp"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const PingHelpMessage = "!ping - check if the bot is running"

// PingCommand command definition
type PingCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ProcessText run command
func (handler *PingCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!ping$`)

	if commandPattern.MatchString(text) {
		return "pong!!", nil
	}

	return "", ErrNextCommand
}
