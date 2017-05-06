package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
)

type StatsCommand struct {
	Db     db.ZbotDatabase
	Next   HandlerCommand
	Levels Levels
}

func (handler *StatsCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!stats$`)
	result := ""

	if commandPattern.MatchString(text) {
		statTotal, err := handler.Db.Statistics()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("Count: %s", statTotal)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
