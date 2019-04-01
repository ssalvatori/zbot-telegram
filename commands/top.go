package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

// TopCommand definition
type TopCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *TopCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!top$`)

	if commandPattern.MatchString(text) {
		//if user.IsAllow(handler.Levels.Top) {
		items, err := handler.Db.Top()
		if err != nil {
			log.Error(err)
			return "", err
		}
		result := fmt.Sprintf(strings.Join(getTerms(items), " "))
		return result, nil
		//} else {
		//	result = fmt.Sprintf("Your level is not enough < %d", handler.Levels.Top)
		//}
	}

	return "", errors.New("text doesn't match")

}
