package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"regexp"
	"strings"
)

type FindCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *FindCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!find\s(\S*)`)
	result := ""

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Find(term[1])
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("%s", strings.Join(getTerms(results), " "))
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
