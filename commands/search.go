package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"

	"log/slog"
)

const SearchHelpMessage = "!search <term> - search for similar <term>"

// SearchCommand definition
type SearchCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

var findSearch = 10

// ProcessText Run module
func (handler *SearchCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!search\s(\S*)`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Search(term[1], uint(channelID), findSearch)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] wasn't found in any definition", term[1]), nil
			}
			slog.Error("search error", "err", err)
			return "", fmt.Errorf("internal error, check logs")
		}
		return strings.Join(getTerms(results), " "), nil
	}
	return "", ErrNextCommand
}
