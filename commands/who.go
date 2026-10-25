package command

import (
	"errors"
	"fmt"
	"regexp"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
)

const WhoHelpMessage = "!who <term> - get information about a term"
const WhoDescription = "Get information about a term"

// WhoCommand definition
type WhoCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// SetDb set db connection if the module need it
func (handler *WhoCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
}

// ShowHelp returns the help message for the command.
func (handler *WhoCommand) ShowHelp(args ...string) string { return WhoHelpMessage }
func (handler *WhoCommand) ShowDescription() string        { return WhoDescription }

// ProcessText run command
func (handler *WhoCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}

		term := commandPattern.FindStringSubmatch(text)
		item := db.Definition{
			Term: term[1],
		}
		Item, err := handler.Db.Get(item.Term, uint(channelID))
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return "", fmt.Errorf("definition [%s] not found", item.Term)
			}
			slog.Error("who get error", "err", err)
			return "", ErrInternalError
		}

		err = handler.Db.IncreaseHits(Item.ID)
		if err != nil {
			slog.Error("who increase hits error", "err", err)
			return "", ErrInternalError
			// if !errors.Is(err, db.ErrInternalError) {
			// 	return "Internal Error, check logs", nil
			// }
		}

		return fmt.Sprintf("[%s] by [%s] on [%v] hits [%d]", Item.Term, Item.Author, utils.ConvertToDateToUTC(Item.UpdatedAt), Item.Hits), nil
	}

	return "", ErrNextCommand
}
