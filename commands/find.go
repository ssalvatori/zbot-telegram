package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"strings"
	log "github.com/Sirupsen/logrus"
)

type FindCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
}

func (handler *FindCommand) ProcessText(text string) string{

	commandPattern := regexp.MustCompile(`^!find\s(\S*)`)
	result := ""

	if(commandPattern.MatchString(text)) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Find(term[1])
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("%s", strings.Join(getTerms(results), " "))
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text)
		}
	}
	return result
}
