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
func (handler *StatsCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!stats$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
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
