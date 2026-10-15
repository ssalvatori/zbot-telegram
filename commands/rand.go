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

const RandHelpMessage = "!rand <limit> <limit> - get a random term, with an optional <limit> (default 1, max 20)"
const RAND_MAX = 20

// RandCommand definition
type RandCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ProcessText run command
func (handler *RandCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!rand(\s+(\d+))?$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		limit := 1
		var err error

		if len(term) == 3 && term[2] != "" {
			limit, err = strconv.Atoi(term[2])

			if err != nil {
				slog.Error("problem converting limit", "term", term[2])
				limit = 1
			}
		}

		if limit > RAND_MAX {
			limit = RAND_MAX
		}

		items, err := handler.Db.Rand(channelName, limit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return "no results", nil
			}
			slog.Error("rand query error", "err", err)
			return "", fmt.Errorf("internal error, check logs")
		}
		var output []string
		for _, item := range items {
			output = append(output, fmt.Sprintf("[%s] - [%s]", item.Term, item.Meaning))

		}

		return strings.Join(output, "\n\n\n"), nil
	}
	return "", ErrNextCommand
}
