package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// LastCommand definition
type LastCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *LastCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *LastCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!last$`)

	if commandPattern.MatchString(text) {
		lastItem, err := handler.Db.Last()
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("[%s] - [%s]", lastItem.Term, lastItem.Meaning), nil
	}
	return "", errors.New("text doesn't match")
}
