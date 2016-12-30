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

	commandPattern := regexp.MustCompile(`^!version$`)

	if(commandPattern.MatchString(text)) {
		return fmt.Sprintf("zbot golang version %s", handler.Version)
	} else {
		return handler.Next.ProcessText(text)
	}

}
