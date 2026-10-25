package command

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"log/slog"

	"github.com/ssalvatori/zbot-telegram/user"
)

// ExecCommand handler to exec.Command
var ExecCommand = exec.Command

// LookPathCommand handler to exec.LookPath
var LookPathCommand = exec.LookPath

const ExternalHelpMessage = "run an external module"
const ExternalDescription = "Run an external module"

// ExternalCommand definition
type ExternalCommand struct {
	PathModules string
}

// ShowHelp returns the help message for the command.
func (handler *ExternalCommand) ShowHelp(args ...string) string { return ExternalHelpMessage }
func (handler *ExternalCommand) ShowDescription() string        { return ExternalDescription }

// ProcessText run command
func (handler *ExternalCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool, isCommand bool) (string, error) {

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

		return handler.RunCommand(fullPathToBinary, user.Username, strconv.Itoa(user.Level), channelName, strings.TrimSpace(args[2])), nil

	}
	return "", nil
}

// RunCommand run external command, the bot is providing the following arguments <username> <level> <channelName> <command_arguments>
func (handler *ExternalCommand) RunCommand(command string, args ...string) string {
	cmd := ExecCommand(command, args...)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("external command error", "output", string(output), "err", err)
		return ""
	}
	return string(output)
}
