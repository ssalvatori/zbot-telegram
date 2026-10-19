package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"log/slog"

	"github.com/go-viper/mapstructure/v2"
)

// GetCurrentDirectory Return the current path
func GetCurrentDirectory() string {
	ex, err := os.Getwd()
	if err != nil {
		slog.Error("could not get working directory", "err", err)
		return os.Getenv("PWD")
	}
	return ex
}

// ConvertToDateToUTC convert unix timestamp to dd-mm-YYYY hh:mm:ss
func ConvertToDateToUTC(unixtime int64) string {
	unixTimeUTC := time.Unix(unixtime, 0).In(time.UTC)
	return fmt.Sprint(unixTimeUTC)
}

// StringToArray parse a string "," as separeator
func StringToArray(cmds string) []string {

	if len(cmds) == 0 {
		return []string{}
	}

	cmdList := strings.Split(cmds, ",")

	for i := range cmdList {
		cmdList[i] = strings.TrimSpace(cmdList[i])
	}

	return cmdList
}

var execCommand = exec.Command

// RunExternalCommand Run external file with a set of arguments
func RunExternalCommand(command string, args ...string) string {
	output, err := execCommand(command, args...).CombinedOutput()
	if err != nil {
		slog.Error("external command failed", "output", string(output), "err", err)
		return ""
	}
	return string(output)
}

// ParseCommand Parse and return command, expecting /command [arg1 arg2]
func ParseCommand(text string) (string, error) {
	slog.Debug("Parsing text", "text", text)

	commandPattern := regexp.MustCompile(`^\/([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)
	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)
		cmd := args[1]
		return cmd, nil
	}
	return "", errors.New("text could not been parser")
}

// GetCommandFile Find file using command key
func GetCommandFile(cmd string, m any) (string, error) {

	type ExternalModule struct {
		Key         string
		File        string
		Description string
	}

	modules := []ExternalModule{}
	err := mapstructure.Decode(m, &modules)
	if err != nil {
		slog.Error("decode error", "err", err)
		return "", fmt.Errorf("command %s not found in list of commands", cmd)
	}

	for i := range modules {
		if modules[i].Key == cmd {
			return modules[i].File, nil
		}
	}
	return "", fmt.Errorf("command %s not found in list of commands", cmd)
}
