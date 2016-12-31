package command

import (
	"regexp"
)

type PingCommand struct {
	Next HandlerCommand
}

func (handler *PingCommand) ProcessText(text string, user User) string {

	commandPattern := regexp.MustCompile(`^!ping$`)
	result := ""

	if(commandPattern.MatchString(text)) {
		result = "pong!!"
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}

