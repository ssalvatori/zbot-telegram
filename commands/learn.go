package command

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//LearnCommand defintion
type LearnCommand struct {
	Db db.ZbotDatabase
}

//ProcessText Run module
func (handler *LearnCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`(?s)^!learn\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
		nowDate := time.Now().Format("2006-01-02")
		term := commandPattern.FindStringSubmatch(text)
		def := db.DefinitionItem{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			Date:    nowDate,
		}
		usedTerm, err := handler.Db.Set(def)
		if err != nil {
			log.Error()
			return "", err
		}
		return fmt.Sprintf("[%s] - [%s]", usedTerm, def.Meaning), nil
	}

	return "", errors.New("text doesn't match")
}
