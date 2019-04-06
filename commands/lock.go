package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

//LockCommand definition
type LockCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *LockCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!lock\s(\S*)$`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		def := db.DefinitionItem{
			Author: user.Username,
			Term:   term[1],
		}
		err := handler.Db.Lock(def)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("[%s] locked", def.Term), nil
	}

	return "", errors.New("text doesn't match")
}
