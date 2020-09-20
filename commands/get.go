package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//GetCommand definition
type GetCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *GetCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^\?\s(\S*)`)
	var result string

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		definition, err := handler.Db.Get(strings.ToLower(term[1]), chat)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] Not found!", term[1]), nil
			}
			log.Error(err.Error())
			return "", ErrInternalError

		}
		err = handler.Db.IncreaseHits(definition.ID)
		if err != nil {
			log.Error(err.Error())
			return "", ErrInternalError
		}
		result = fmt.Sprintf("[%s] - [%s]", definition.Term, definition.Meaning)
		//		}
		return result, nil
	}
	return "", ErrNextCommand
}
