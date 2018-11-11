package command

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type LearnCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

//ProcessText Run module
func (handler *LearnCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!learn\s(\S*)\s(.*)`)
	result := ""

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
			log.Error(fmt.Errorf("Error learn %v", err))
			return ""
		}
		result = fmt.Sprintf("[%s] - [%s]", usedTerm, def.Meaning)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
