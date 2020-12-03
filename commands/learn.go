package command

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//LearnCommand defintion
type LearnCommand struct {
	Db db.ZbotDatabase
}

//ProcessText Run module
func (handler *LearnCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`(?s)^!learn\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
		// nowDate := time.Now().Format("2006-01-02")
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			// Date:    nowDate,
			Chat: chat,
		}
		usedTerm, err := handler.Db.Set(def)
		if err != nil {
			log.Error()
			return "", err
		}
		return fmt.Sprintf("[%s] - [%s]", usedTerm, def.Meaning), nil
	}

	return "", ErrNextCommand
}
