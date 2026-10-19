package command

import (
	"fmt"
	"regexp"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const LearnHelpMessage = "!learn <term> <meaning> - learn a <term> with <meaning>; duplicate terms receive a numeric suffix"

// LearnCommand defintion
type LearnCommand struct {
	Db db.ZbotDatabase
}

// ProcessText Run module
func (handler *LearnCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`(?s)^!learn\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
		// nowDate := time.Now().Format("2006-01-02")
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			// Date:    nowDate,
			ChannelID: uint(channelID),
		}
		usedTerm, err := handler.Db.Set(def)
		if err != nil {
			slog.Error("learn db error", "err", err)
			return "", err
		}
		return fmt.Sprintf("[%s] - [%s]", usedTerm, def.Meaning), nil
	}

	return "", ErrNextCommand
}
