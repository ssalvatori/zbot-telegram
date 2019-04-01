package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// StatsCommand definition
type StatsCommand struct {
	Db db.ZbotDatabase
	//	Next   HandlerCommand
	//	Levels Levels
}

// ProcessText run command
func (handler *StatsCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!stats$`)

	if commandPattern.MatchString(text) {
		//if user.IsAllow(handler.Levels.Stats) {
		statTotal, err := handler.Db.Statistics()
		if err != nil {
			log.Error(err)
			return "", err
		}
		result := fmt.Sprintf("Count: %s", statTotal)
		//} else {
		//	result = fmt.Sprintf("Your level is not enough < %d", handler.Levels.Stats)
		//}
		return result, nil

	}
	//	return "", result
	return "", errors.New("text doesn't match")
}
