package command

import (
	"regexp"
	"os/exec"
	"bytes"
	log "github.com/Sirupsen/logrus"
)

type ExternalCommand struct {
	PathModules string
	Next HandlerCommand
}

func (handler *ExternalCommand) ProcessText(text string, user User) string {


	commandPattern := regexp.MustCompile(`^!([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)
	result := ""

	if(commandPattern.MatchString(text)) {
		args := commandPattern.FindStringSubmatch(text)
		externalModule := args[1]

		binary, err := exec.LookPath(handler.PathModules+externalModule)

		if err != nil {
			log.Error(err)
			return ""
		}

		cmd := exec.Command(binary, args[2])
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error(err)
			return ""
		}
		result = out.String()
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}


