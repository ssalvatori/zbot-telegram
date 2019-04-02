package command

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// AppendCommand definition
type AppendCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *AppendCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`(?s)^!append\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
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
			return "", err
		}
		def, err = handler.Db.Get(def.Term)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("[%s] = [%s]", def.Term, def.Meaning), nil
	}

	return "", errors.New("text doesn't match")
}
