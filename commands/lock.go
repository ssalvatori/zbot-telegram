package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"regexp"
)

type LockCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *LockCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!lock\s(\S*)$`)
	result := ""

	if commandPattern.MatchString(text) {
		if IsUserAllow(handler.Db, user.Username, handler.Levels.Lock) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Author: user.Username,
				Term:   term[0],
			}
			err := handler.Db.Lock(def)
			if err != nil {
				log.Error(err)
			}

		} else {
			result = fmt.Sprintf("Your level is not enough < %s", handler.Levels.Lock)
		}
	}

	return result
}
