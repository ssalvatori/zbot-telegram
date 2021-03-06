package command

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/user"
)

//ExecCommand handler to exec.Command
var ExecCommand = exec.Command

//LookPathCommand handler to exec.LookPath
var LookPathCommand = exec.LookPath

//ExternalCommand definition
type ExternalCommand struct {
	PathModules string
}

// ProcessText run command
func (handler *ExternalCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)

	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)
		externalModule := args[1]

		log.Debug("Looking for module: " + handler.PathModules + externalModule)

		fullPathToBinary, err := LookPathCommand(handler.PathModules + externalModule)

		if err != nil {
			log.Error(err)
			// return "", fmt.Errorf("Internal error with command [%s]", externalModule)
			return "", nil
		}

		return handler.RunCommand(fullPathToBinary, user.Username, strconv.Itoa(user.Level), chat, strings.TrimSpace(args[2])), nil

	}
	return "", nil
}

//RunCommand run external command, the bot is providing the following arguments <username> <level> <chat> <command_arguments>
func (handler *ExternalCommand) RunCommand(command string, args ...string) string {
	output, err := ExecCommand(command, args...).CombinedOutput()
	if err != nil {
		log.Error(fmt.Sprintf("%s", output))
		log.Error(err)
		return ""
	}
	return fmt.Sprintf("%s", output)
}
