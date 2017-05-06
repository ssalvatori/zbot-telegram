package command

import (
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
)

type PingCommand struct {
	Next   HandlerCommand
	Levels Levels
}

func (handler *PingCommand) ProcessText(text string, user user.User) string {

	commandPattern := regexp.MustCompile(`^!ping$`)
	result := ""

	if commandPattern.MatchString(text) {
		result = "pong!!"
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
