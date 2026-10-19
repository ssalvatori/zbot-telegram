package command

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const TopHelpMessage = "!top <limit> - get the top terms by hits"
const TOP_MAX = 20

// TopCommand definition
type TopCommand struct {
	Db db.ZbotDatabase
}

// SetDb set db connection if the module need it
func (handler *TopCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
}

// ProcessText run command
func (handler *TopCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!top(\s+(\d+))?$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		limit := TOP_MAX
		var err error

		if len(term) == 3 && term[2] != "" {
			limit, err = strconv.Atoi(term[2])

			if err != nil {
				slog.Error("problem converting limit", "term", term[2])
				limit = TOP_MAX
			}
		}

		if limit > TOP_MAX {
			limit = TOP_MAX
		}

		items, err := handler.Db.Top(uint(channelID), limit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return "no results", nil
			}
			slog.Error("top query error", "err", err)
			return "", fmt.Errorf("internal error, check logs")
		}
		return strings.Join(getTerms(items), " "), nil
	}

	return "", ErrNextCommand

}
