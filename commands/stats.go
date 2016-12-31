package command

import (
	"regexp"
	"github.com/ssalvatori/zbot-telegram-go/db"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

type StatsCommand struct {
	Db db.ZbotDatabase
	Next     HandlerCommand
}

func (handler *StatsCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!stats$`)
	result := ""

	if(commandPattern.MatchString(text)) {
		statTotal, err := handler.Db.Statistics()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("Count: %s",statTotal)
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}

