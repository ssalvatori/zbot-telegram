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
	result := ""

	if(commandPattern.MatchString(text)) {
		result = fmt.Sprintf("zbot golang version %s", handler.Version)
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text)
		}
	}
	return result

}
