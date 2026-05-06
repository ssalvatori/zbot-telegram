package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"log/slog"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//GetCommand definition
type GetCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *GetCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^\?\s(\S*)`)
	var result string

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		definition, err := handler.Db.Get(strings.ToLower(term[1]), chat)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] Not found!", term[1]), nil
			}
			slog.Error("get error", "err", err)
			return "", ErrInternalError

		}
		err = handler.Db.IncreaseHits(definition.ID)
		if err != nil {
			slog.Error("increase hits error", "err", err)
			return "", ErrInternalError
		}
		result = fmt.Sprintf("[%s] - [%s]", definition.Term, definition.Meaning)
		//		}
		return result, nil
	}
	return "", ErrNextCommand
}
