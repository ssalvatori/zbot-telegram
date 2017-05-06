package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
)

type LockCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *LockCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!lock\s(\S*)$`)
	result := ""

	if commandPattern.MatchString(text) {
		if user.IsAllow(handler.Levels.Lock) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Author: user.Username,
				Term:   term[1],
			}
			err := handler.Db.Lock(def)
			if err != nil {
				log.Error(err)
			}
			result = fmt.Sprintf("[%s] locked", def.Term)

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
