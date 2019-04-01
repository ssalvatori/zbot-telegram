package command

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram-go/user"
)

//VersionCommand configuration for version
type VersionCommand struct {
	//Next      HandlerCommand
	Version   string
	GitHash   string
	BuildTime string
}

//ProcessText run command
func (handler *VersionCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!version$`)

	if commandPattern.MatchString(text) {
		return fmt.Sprintf("zbot golang version [%s] commit [%s] build-time [%s]", handler.Version, handler.GitHash, handler.BuildTime), nil
	}
	return "", errors.New("text doesn't match")

}
