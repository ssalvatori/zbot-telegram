package command

import (
	"regexp"
	"fmt"
)

type VersionCommand struct {
	Next HandlerCommand
	Version string
}

func (handler *VersionCommand) ProcessText(text string) string{

	var command string = "^!version"

	if(regexp.MatchString(regexp.MustCompile(command), text)) {
		return fmt.Sprintf("zbot golang version %s", handler.version)
	} else {
		return handler.next.Process(text)
	}

}
