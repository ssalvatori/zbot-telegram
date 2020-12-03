package command

import (
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// ForgetCommand definition
type ForgetCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *ForgetCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!forget\s(\S*)$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term: term[1],
		}
		err := handler.Db.Forget(def, chat)
		if err != nil {
			return "", ErrInternalError
		}
		return fmt.Sprintf("[%s] deleted", def.Term), nil
	}
	return "", ErrNextCommand
}
