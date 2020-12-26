package command

import (
	"regexp"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// PingCommand command definition
type PingCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *PingCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!ping$`)

	if commandPattern.MatchString(text) {
		return "pong!!", nil
	}

	return "", ErrNextCommand
}
