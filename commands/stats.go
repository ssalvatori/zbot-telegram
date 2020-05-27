package command

import (
	"errors"
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
func (handler *StatsCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!stats$`)

	if commandPattern.MatchString(text) {
		statTotal, err := handler.Db.Statistics()
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("Count: %s", statTotal), nil

	}
	//	return "", result
	return "", errors.New("text doesn't match")
}
