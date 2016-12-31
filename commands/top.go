package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	log "github.com/Sirupsen/logrus"
	"strings"
)

type TopCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
}

func (handler *TopCommand) ProcessText(text string, user User) string{

	commandPattern := regexp.MustCompile(`^!top$`)
	result := ""

	if(commandPattern.MatchString(text)) {
		items, err := handler.Db.Top()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf(strings.Join(getTerms(items), " "))
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}