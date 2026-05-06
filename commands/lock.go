package command

import (
	"fmt"
	"regexp"

	"log/slog"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//LockCommand definition
type LockCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *LockCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!lock\s(\S*)$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Author: user.Username,
			Term:   term[1],
		}
		err := handler.Db.Lock(def, chat)
		if err != nil {
			slog.Error("lock error", "err", err)
			return "", err
		}
		return fmt.Sprintf("[%s] locked", def.Term), nil
	}

	return "", ErrNextCommand
}
