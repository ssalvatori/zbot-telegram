package command

import (
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

const ForgetHelpMessage = "!forget <term> - delete a <term> from the dictionary"

// ForgetCommand definition
type ForgetCommand struct {
	Db          db.ZbotDatabase
	HelpMessage string
}

// ProcessText run command
func (handler *ForgetCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!forget\s(\S*)$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(channelName) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term:      term[1],
			ChannelID: uint(channelID),
		}
		err := handler.Db.Forget(def)
		if err != nil {
			return "", ErrInternalError
		}
		return fmt.Sprintf("[%s] deleted", def.Term), nil
	}
	return "", ErrNextCommand
}
