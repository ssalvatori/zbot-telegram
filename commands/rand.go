package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	log "github.com/Sirupsen/logrus"
)

type RandCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
}

func (handler *RandCommand) ProcessText(text string) string{

	commandPattern := regexp.MustCompile(`^!rand$`)
	result := ""

	if(commandPattern.MatchString(text)) {
		randItem, err := handler.Db.Rand()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("[%s] - [%s]", randItem.Term, randItem.Meaning)
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text)
		}
	}
	return result
}
