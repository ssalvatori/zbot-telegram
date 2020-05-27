package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// TopCommand definition
type TopCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *TopCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *TopCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!top$`)

	if commandPattern.MatchString(text) {
		items, err := handler.Db.Top()
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf(strings.Join(getTerms(items), " ")), nil
	}

	return "", errors.New("text doesn't match")

}
