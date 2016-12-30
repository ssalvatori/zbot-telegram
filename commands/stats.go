package command

import (
	"regexp"
	"github.com/ssalvatori/zbot-telegram-go/database"
)

type StatsCommand struct {
	Db zbotDatabase
	Next     HandlerCommand
}

func (handler *StatsCommand) ProcessText(text string) string {

	commandPattern := regexp.MustCompile(`^!stats`)

	if(commandPattern.MatchString(text)) {
		return Db.statistics()
	} else {
		return handler.Next.ProcessText(text)
	}

}

