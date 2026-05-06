package command

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
)

// ExecCommand handler to exec.Command
var ExecCommand = exec.Command

// LookPathCommand handler to exec.LookPath
var LookPathCommand = exec.LookPath

// ExternalCommand definition
type ExternalCommand struct {
	PathModules string
}

// ProcessText run command
func (handler *ExternalCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	commandPattern := regexp.MustCompile(`^!([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)

	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)
		externalModule := args[1]

		slog.Debug("Looking for module", "path", handler.PathModules+externalModule)

		fullPathToBinary, err := LookPathCommand(handler.PathModules + externalModule)

		if err != nil {
			slog.Error("module lookup error", "err", err)
			// return "", fmt.Errorf("Internal error with command [%s]", externalModule)
			return "", nil
		}

		return handler.RunCommand(fullPathToBinary, utils.SanitizeArg(user.Username), strconv.Itoa(user.Level), utils.SanitizeArg(chat), utils.SanitizeArg(strings.TrimSpace(args[2]))), nil

	}
	return "", nil
}

// RunCommand run external command, the bot is providing the following arguments <username> <level> <chat> <command_arguments>
func (handler *ExternalCommand) RunCommand(command string, args ...string) string {
	cmd := ExecCommand(command, args...)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("external command error", "output", string(output), "err", err)
		return ""
	}
	return fmt.Sprintf("%s", output)
}
