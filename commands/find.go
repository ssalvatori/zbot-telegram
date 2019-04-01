package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

//FindCommand defintion
type FindCommand struct {
	Db db.ZbotDatabase
}

//ProcessText run command
func (handler *FindCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!find\s(\S*)`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Find(term[1])
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("%s", strings.Join(getTerms(results), " ")), nil
	}
	return "", errors.New("text doesn't match")
}
