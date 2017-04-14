package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"regexp"
)

type ForgetCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *ForgetCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!forget\s(\S*)$`)
	result := ""

	if commandPattern.MatchString(text) {
		if IsUserAllow(handler.Db, user.Username, handler.Levels.Forget) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Term: term[1],
			}
			err := handler.Db.Forget(def)
			if err != nil {
				log.Error(err)
			}
			return fmt.Sprintf("[%s] deleted", def.Term)
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}

	return result
}
