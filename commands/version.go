package command

import (
	"fmt"
	"regexp"

	"github.com/ssalvatori/zbot-telegram-go/user"
)

//VersionCommand configuration for version
type VersionCommand struct {
	Next      HandlerCommand
	Version   string
	GitHash   string
	BuildTime string
}

//ProcessText run command
func (handler *VersionCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!version$`)
	result := ""

	if commandPattern.MatchString(text) {
		result = fmt.Sprintf("zbot golang version [%s] commit [%s] build-time [%s]", handler.Version, handler.GitHash, handler.BuildTime)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result

}
