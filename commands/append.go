package command

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type AppendCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *AppendCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!append\s(\S*)\s(.*)`)
	result := ""

	if commandPattern.MatchString(text) {
		if user.IsAllow(handler.Levels.Append) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Term:    term[1],
				Meaning: term[2],
				Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
				Date:    time.Now().Format("2006-01-02"),
			}
			err := handler.Db.Append(def)
			if err != nil {
				log.Error(err)
			}
			def, err = handler.Db.Get(def.Term)
			if err != nil {
				log.Error(err)
			}
			result = fmt.Sprintf("[%s] = [%s]", def.Term, def.Meaning)
		} else {
			result = fmt.Sprintf("Your level is not enough < %s", handler.Levels.Lock)
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}

	return result
}
