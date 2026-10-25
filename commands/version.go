package command

import (
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram/user"
)

const VersionHelpMessage = "Show bot version"
const VersionDescription = "Show the bot version, git hash and build time"

// VersionCommand configuration for version
type VersionCommand struct {
	Version   string
	GitHash   string
	BuildTime string
}

// ShowHelp returns the help message for the command.
func (handler *VersionCommand) ShowHelp(args ...string) string { return VersionHelpMessage }
func (handler *VersionCommand) ShowDescription() string        { return VersionDescription }

// ProcessText run command
func (handler *VersionCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!version$`)

	if isCommand || commandPattern.MatchString(text) {
		return fmt.Sprintf("zbot golang version [%s] commit [%s] build-time [%s]", handler.Version, handler.GitHash, handler.BuildTime), nil
	}
	return "", ErrNextCommand

}
