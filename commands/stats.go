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
}

//SetDb set db connection if the module need it
func (handler *StatsCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
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
