package command

import (
	"errors"
	"regexp"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// PingCommand command definition
type PingCommand struct {
}

//SetDb set db connection if the module need it
func (handler *PingCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *PingCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!ping$`)

	if commandPattern.MatchString(text) {
		return "pong!!", nil
	}

	return "", errors.New("text doesn't match")
}
