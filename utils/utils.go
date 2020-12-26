package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// InArray returns true if the string 's' is found in the array 'arr', otherwise false
func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

//GetCurrentDirectory Return the current path
func GetCurrentDirectory() string {
	ex, err := os.Getwd()
	if err != nil {
		log.Error(fmt.Errorf("Could get the path %v", err))
		return os.Getenv("PWD")
	}
	return ex
}

//ConvertToDateToUTC convert unix timestamp to dd-mm-YYYY hh:mm:ss
func ConvertToDateToUTC(unixtime int64) string {
	location, _ := time.LoadLocation("UTC")
	unixTimeUTC := time.Unix(unixtime, 0).In(location)
	return fmt.Sprint(unixTimeUTC)
}

//StringToArray parse a string "," as separeator
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

//RunExternalCommand Run external file with a set of arguments
func RunExternalCommand(command string, args ...string) string {
	output, err := execCommand(command, args...).CombinedOutput()
	if err != nil {
		log.Error(fmt.Sprintf("%s", output))
		log.Error(err)
		return ""
	}
	return fmt.Sprintf("%s", output)
}

//ParseCommand Parse and return command, expecting /command [arg1 arg2]
func ParseCommand(text string) (string, error) {
	log.Debug(fmt.Sprintf("Parsing text %s", text))

	commandPattern := regexp.MustCompile(`^\/([a-zA-Z0-9\_\-]+)([\s(\S*)]*)?`)
	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)
		cmd := args[1]
		return cmd, nil
	}
	return "", errors.New("Text could not been parser")
}

//GetCommandFile Find file using command key
func GetCommandFile(cmd string, modules []struct {
	Key         string
	File        string
	Description string
}) (string, error) {
	for i := range modules {
		if modules[i].Key == cmd {
			return modules[i].File, nil
		}
	}
	return "", fmt.Errorf("Command %s not found in list of commands", cmd)
}
