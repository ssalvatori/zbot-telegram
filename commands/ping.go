package command

import (
	"regexp"
)

type PingCommand struct {
	Next HandlerCommand
}

func (handler *PingCommand) ProcessText(text string) string {

	commandPattern := regexp.MustCompile(`^!ping`)

	if(commandPattern.MatchString(text)) {
		return "pong!!"
	} else {
		return handler.Next.ProcessText(text)
	}

}

