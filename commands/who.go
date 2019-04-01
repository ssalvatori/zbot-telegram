package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// WhoCommand definition
type WhoCommand struct {
	//Next   HandlerCommand
	Db db.ZbotDatabase
	//Levels Levels
}

//ProcessText run command
func (handler *WhoCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)

	if commandPattern.MatchString(text) {
		//if user.IsAllow(handler.Levels.Who) {
		term := commandPattern.FindStringSubmatch(text)
		def := db.DefinitionItem{
			Term: term[1],
		}
		Item, err := handler.Db.Get(def.Term)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		result := fmt.Sprintf("[%s] by [%s] on [%s]", Item.Term, Item.Author, Item.Date)
		return result, nil
		//}
	}
	/*else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}*/

	return "", errors.New("text doesn't match")
}
