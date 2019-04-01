package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

//RandCommand definition
type RandCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *RandCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
}

// ProcessText run command
func (handler *RandCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!rand$`)

	if commandPattern.MatchString(text) {
		randItem, err := handler.Db.Rand()
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("[%s] - [%s]", randItem.Term, randItem.Meaning), nil
	}
	return "", errors.New("text doesn't match")
}
