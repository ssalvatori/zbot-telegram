package command

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

//ExternalCommand definition
type ExternalCommand struct {
	PathModules string
}

// ProcessText run command
func (handler *ExternalCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)

	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)
		externalModule := args[1]

		log.Debug("Looking for module: " + handler.PathModules + externalModule)

		binary, err := exec.LookPath(handler.PathModules + externalModule)

		if err != nil {
			log.Error(err)
			return "", err
		}

		cmd := exec.Command(binary, user.Username, strconv.Itoa(user.Level), args[2])
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Error(err)
			return "", err
		}
		return out.String(), nil
	}
	return "", nil
}
