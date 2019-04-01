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
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *WhoCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
}

//ProcessText run command
func (handler *WhoCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)

	if commandPattern.MatchString(text) {

		term := commandPattern.FindStringSubmatch(text)
		def := db.DefinitionItem{
			Term: term[1],
		}
		Item, err := handler.Db.Get(def.Term)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}

		return fmt.Sprintf("[%s] by [%s] on [%s]", Item.Term, Item.Author, Item.Date), nil

	}

	return "", errors.New("text doesn't match")
}
