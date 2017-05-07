package command

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type GetCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *GetCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^\?\s(\S*)`)
	result := ""

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		definition, err := handler.Db.Get(strings.ToLower(term[1]))
		if err != nil {
			log.Error(err)
		}
		if definition.Term != "" {
			result = fmt.Sprintf("[%s] - [%s]", definition.Term, definition.Meaning)
		} else {
			result = fmt.Sprintf("[%s] Not found!", term[1])
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
