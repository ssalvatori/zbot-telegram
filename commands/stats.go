package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
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
		if user.IsAllow(handler.Levels.Stats) {
			statTotal, err := handler.Db.Statistics()
			if err != nil {
				log.Error(err)
				return ""
			}
			result = fmt.Sprintf("Count: %s", statTotal)
		} else {
			result = fmt.Sprintf("Your level is not enough < %d", handler.Levels.Stats)
		}

	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
