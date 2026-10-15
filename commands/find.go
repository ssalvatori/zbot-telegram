package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const FindHelpMessage = "!find <text> - find _terms_ whose meaning contains <text>."

// FindCommand definition
type FindCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

var findLimit = 10

// ProcessText run command
func (handler *FindCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!find\s(\S*)`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Find(term[1], channelName, findLimit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] wasn't found in the content of any definition", term[1]), nil
			}
			return "", ErrInternalError
		}
		return strings.Join(getTerms(results), " "), nil
	}
	return "", ErrNextCommand
}
