package command

import (
	"fmt"
	"regexp"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const StatsHelpMessage = "!stats - get the total number of terms"

// StatsCommand definition
type StatsCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ProcessText run command
func (handler *StatsCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!stats$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		statTotal, err := handler.Db.Statistics(uint(channelID))
		if err != nil {
			slog.Error("statistics error", "err", err)
			return "", db.ErrInternalError
		}
		return fmt.Sprintf("Number of definitions: %s", statTotal), nil

	}
	//	return "", result
	return "", ErrNextCommand
}
