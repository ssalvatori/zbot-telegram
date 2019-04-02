package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

//GetCommand definition
type GetCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *GetCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^\?\s(\S*)`)
	var result string

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		definition, err := handler.Db.Get(strings.ToLower(term[1]))
		if err != nil {
			log.Error(err)
			return "", err
		}
		if definition.Term != "" {
			result = fmt.Sprintf("[%s] - [%s]", definition.Term, definition.Meaning)
		} else {
			result = fmt.Sprintf("[%s] Not found!", term[1])
		}
		return result, nil
	}
	return "", errors.New("text doesn't match")
}
