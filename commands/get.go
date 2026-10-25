package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const GetHelpMessage = "!get <term> - get a <term>'s meaning"
const GetDescription = "Get the meaning of a term."

// GetCommand definition
type GetCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ShowHelp returns the help message for the command.
func (handler *GetCommand) ShowHelp(args ...string) string { return GetHelpMessage }
func (handler *GetCommand) ShowDescription() string        { return GetDescription }

// ProcessText run command
func (handler *GetCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^\?\s(\S*)`)
	var result string

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		definition, err := handler.Db.Get(strings.ToLower(term[1]), uint(channelID))
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("[%s] Not found!", term[1]), nil
			}
			slog.Error("get error", "err", err)
			return "", ErrInternalError

		}
		err = handler.Db.IncreaseHits(definition.ID)
		if err != nil {
			slog.Error("increase hits error", "err", err)
			return "", ErrInternalError
		}
		result = fmt.Sprintf("[%s] - [%s]", definition.Term, definition.Meaning)
		//		}
		return result, nil
	}
	return "", ErrNextCommand
}
