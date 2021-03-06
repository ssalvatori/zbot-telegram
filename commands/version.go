package command

import (
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram/user"
)

//VersionCommand configuration for version
type VersionCommand struct {
	Version   string
	GitHash   string
	BuildTime string
}

//ProcessText run command
func (handler *VersionCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!version$`)

	if commandPattern.MatchString(text) {
		return fmt.Sprintf("zbot golang version [%s] commit [%s] build-time [%s]", handler.Version, handler.GitHash, handler.BuildTime), nil
	}
	return "", ErrNextCommand

}
