package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"regexp"
)

type WhoCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *WhoCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)
	result := ""

	if commandPattern.MatchString(text) {
		if IsUserAllow(handler.Db, user.Username, handler.Levels.Append) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Term: term[1],
			}
			Item, err := handler.Db.Get(def.Term)
			if err != nil {
				log.Error(err)
			}
			result = fmt.Sprintf("[%s] by [%s] on [%s]", Item.Term, Item.Author, Item.Date)
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}

	return result
}
