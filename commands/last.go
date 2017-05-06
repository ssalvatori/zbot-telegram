package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type LastCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *LastCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!last$`)
	result := ""

	if commandPattern.MatchString(text) {
		lastItem, err := handler.Db.Last()
		if err != nil {
			log.Error(err)
			return ""
		}
		result = fmt.Sprintf("[%s] - [%s]", lastItem.Term, lastItem.Meaning)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
