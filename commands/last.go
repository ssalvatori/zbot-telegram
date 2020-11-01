package command

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// LastCommand definition
type LastCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *LastCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *LastCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!last$`)

	if commandPattern.MatchString(text) {
		lastItems, err := handler.Db.Last(chat, 10)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return PrintTerms(lastItems), nil
	}
	return "", ErrNextCommand
}

//PrintTerms .
func PrintTerms(items []db.Definition) string {
	keys := make([]string, 0, len(items))
	for item := range items {
		keys = append(keys, items[item].Term)
	}
	return fmt.Sprintf("[ %s ]", strings.Join(keys, " "))
}
