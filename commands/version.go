package command

import (
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
)

type VersionCommand struct {
	Next      HandlerCommand
	Version   string
	BuildTime string
	Levels    Levels
}

func (handler *VersionCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!version$`)
	result := ""

	if commandPattern.MatchString(text) {
		result = fmt.Sprintf("zbot golang version [%s] build-time [%s]", handler.Version, handler.BuildTime)
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result

}
