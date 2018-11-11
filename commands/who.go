package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type WhoCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *WhoCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)
	result := ""

	if commandPattern.MatchString(text) {
		if user.IsAllow(handler.Levels.Who) {
			term := commandPattern.FindStringSubmatch(text)
			def := db.DefinitionItem{
				Term: term[1],
			}
			Item, err := handler.Db.Get(def.Term)
			if err != nil {
				log.Error(fmt.Errorf("Error learn %v", err))
				return ""
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
