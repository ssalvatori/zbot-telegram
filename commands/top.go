package command

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
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
		if user.IsAllow(handler.Levels.Top) {
			items, err := handler.Db.Top()
			if err != nil {
				log.Error(err)
				return ""
			}
			result = fmt.Sprintf(strings.Join(getTerms(items), " "))
		} else {
			result = fmt.Sprintf("Your level is not enough < %d", handler.Levels.Top)
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
