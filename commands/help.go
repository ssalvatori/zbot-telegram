package command

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/ssalvatori/zbot-telegram/user"
)

const HelpMessage = "!help <command> - get help for a specific <command>."
const HelpDescription = "Get help for a specific command."

type HelpCommand struct {
	Commands *CommandsList
}

// ShowHelp returns the help message for the command.
func (handler *HelpCommand) ShowHelp(args ...string) string { return HelpMessage }
func (handler *HelpCommand) ShowDescription() string        { return HelpDescription }

func (handler *HelpCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!help\s+(\S+)\s*$`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)

		slog.Debug("help command term", "term", term[1])
		if handler.Commands == nil || handler.Commands.List == nil {
			return "", ErrNextCommand
		}

		for element := handler.Commands.List.Front(); element != nil; element = element.Next() {
			command := element.Value.(*CommandElement)
			if strings.EqualFold(command.CmdString, term[1]) {
				return "", nil
			}
		}

	}

	return "", ErrNextCommand

}
