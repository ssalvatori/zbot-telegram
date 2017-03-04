package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"regexp"
	"time"
)

type LearnCommand struct {
	Next HandlerCommand
	Db   db.ZbotDatabase
}

func (handler *LearnCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!learn\s(\S*)\s(.*)`)
	result := ""

	if commandPattern.MatchString(text) {
		nowDate := time.Now().Format("2006-01-02")
		term := commandPattern.FindStringSubmatch(text)
		def := db.DefinitionItem{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			Date:    nowDate,
		}
		err := handler.Db.Set(def)
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("[%s] - [%s]", def.Term, def.Meaning)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
