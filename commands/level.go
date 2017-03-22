package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

type LevelCommand struct {
	Db     db.ZbotDatabase
	Next   HandlerCommand
	Levels Levels
}

func (handler *LevelCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!level$`)
	result := ""

	if commandPattern.MatchString(text) {
		level, err := handler.Db.UserLevel(user.Username)
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("%s level %s", user.Username, level)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
