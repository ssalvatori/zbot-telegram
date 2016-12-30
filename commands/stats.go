package command

import (
	"regexp"
	"fmt"
)

type StatsCommand struct {
	next HandlerCommand
	version string
}

func (handler *StatsCommand) Process(text string) string {

	var command string = "^!stats"

	if(regexp.MatchString(regexp.MustCompile(command), text)) {
		return fmt.Sprintf("zbot golang version %s", handler.version)
	} else {
		return handler.next.Process(text)
	}

}

