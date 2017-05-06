package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
	"strings"
)

type TopCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

func (handler *TopCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!top$`)
	result := ""

	if commandPattern.MatchString(text) {
		items, err := handler.Db.Top()
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf(strings.Join(getTerms(items), " "))
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
