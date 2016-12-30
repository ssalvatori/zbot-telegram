package command

import (
	"regexp"
	"fmt"
)

type PingCommand struct {
	Next HandlerCommand
	Version string
}

func (handler *PingCommand) ProcessText(text string) string {

	var command string = "^!ping"

	if(regexp.MatchString(regexp.MustCompile(command), text)) {
		return fmt.Sprintf("!pong")
	} else {
		return handler.next.Process(text)
	}

}

