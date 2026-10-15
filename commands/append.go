package command

import (
	"fmt"
	"regexp"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const AppendHelpMessage = "!append <term> <meaning> - append a <meaning> to a <term>."

// AppendCommand definition
type AppendCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ProcessText run command
func (handler *AppendCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`(?s)^!append\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			Chat:    channelName,
		}
		err := handler.Db.Append(def, channelName)
		if err != nil {
			slog.Error("append error", "err", err)
			return "", ErrInternalError
		}

		def, err = handler.Db.Get(def.Term, channelName)
		if err != nil {
			slog.Error("get after append error", "err", err)
			return "", ErrInternalError
		}
		return fmt.Sprintf("[%s] = [%s]", def.Term, def.Meaning), nil
	}

	return "", ErrNextCommand
}
