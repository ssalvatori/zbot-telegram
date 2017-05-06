package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
)

type RandCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *RandCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!rand$`)
	result := ""

	if commandPattern.MatchString(text) {
		randItem, err := handler.Db.Rand()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("[%s] - [%s]", randItem.Term, randItem.Meaning)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
