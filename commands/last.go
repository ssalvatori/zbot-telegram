package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	log "github.com/Sirupsen/logrus"
)

type LastCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
}

func (handler *LastCommand) ProcessText(text string, user User) string{

	commandPattern := regexp.MustCompile(`^!last$`)
	result := ""

	if(commandPattern.MatchString(text)) {
		lastItem, err := handler.Db.Last()
		if err != nil {
			log.Error(err)
		}
		result =fmt.Sprintf("[%s] - [%s]", lastItem.Term, lastItem.Meaning)
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}