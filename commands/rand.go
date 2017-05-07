package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type RandCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}
// ProcessText
func (handler *RandCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!rand$`)
	result := ""

	if commandPattern.MatchString(text) {
		randItem, err := handler.Db.Rand()
		if err != nil {
			log.Error(err)
			return ""
		}
		result = fmt.Sprintf("[%s] - [%s]", randItem.Term, randItem.Meaning)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
