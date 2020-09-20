package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//FindCommand defintion
type FindCommand struct {
	Db db.ZbotDatabase
}

var findLimit = 10

//ProcessText run command
func (handler *FindCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!find\s(\S*)`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Find(term[1], chat, findLimit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] wasn't found in the content of any definition", term[1]), nil
			}
			return "", ErrInternalError
		}
		return fmt.Sprintf("%s", strings.Join(getTerms(results), " ")), nil
	}
	return "", ErrNextCommand
}
