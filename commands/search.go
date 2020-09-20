package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//SearchCommand definition
type SearchCommand struct {
	Db db.ZbotDatabase
}

var findSearch = 10

//ProcessText Run module
func (handler *SearchCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!search\s(\S*)`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Search(term[1], chat, findSearch)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] wasn't found in any definition", term[1]), nil
			}
			if !errors.Is(err, db.ErrInternalError) {
				return "", fmt.Errorf("Internal error, check logs")
			}
		}
		return fmt.Sprintf("%s", strings.Join(getTerms(results), " ")), nil
	}
	return "", ErrNextCommand
}
