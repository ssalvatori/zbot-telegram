package command

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// StatsCommand definition
type StatsCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *StatsCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!stats$`)

	if commandPattern.MatchString(text) {
		statTotal, err := handler.Db.Statistics(chat)
		if err != nil {
			log.Error(err)
			return "", db.ErrInternalError
		}
		return fmt.Sprintf("Number of definitions: %s", statTotal), nil

	}
	//	return "", result
	return "", ErrNextCommand
}
