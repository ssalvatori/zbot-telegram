package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var externalCommand = ExternalCommand{
	PathModules: "",
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeLookPathCommand(file string) (string, error) {
	if file == "mock-error" {
		return "", errors.New("mock")
	}
	return fmt.Sprintf("/home/ssalvatori/module/%s", file), nil
}

func TestRunCommand(t *testing.T) {

	ExecCommand = fakeExecCommand
	LookPathCommand = fakeLookPathCommand

	defer func() {
		ExecCommand = exec.Command
		LookPathCommand = exec.LookPath
	}()

	out := externalCommand.RunCommand("external_module", "ssalvatori", "100")
	assert.Equal(t, "external_module ssalvatori 100\n", out, "Run Command")
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	//fmt.Fprintf(os.Stdout, externalFakeOutput)
	//os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No args\n")
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "%s\n", strings.Join(args, " "))
	os.Exit(0)
}

func TestExternalCommandOK(t *testing.T) {
	userTest.Level = 100

	LookPathCommand = fakeLookPathCommand
	ExecCommand = fakeExecCommand

	defer func() {
		ExecCommand = exec.Command
		LookPathCommand = exec.LookPath
	}()

	result, _ := externalCommand.ProcessText("!external_module arg1 arg2 arg3", userTest)
	assert.Equal(t, "/home/ssalvatori/module/external_module ssalvatori 100  arg1 arg2 arg3\n", result, "external")

	_, err := externalCommand.ProcessText("!mock-error arg1 arg2", userTest)
	assert.Equal(t, "mock", err.Error(), "external mock")

}

func TestExternalCommandInject(t *testing.T) {
	userTest.Level = 100

	result, _ := externalCommand.ProcessText("!../../test arg1 arg2 arg3", userTest)

	assert.Equal(t, "", result, "external commmand inject")
}
