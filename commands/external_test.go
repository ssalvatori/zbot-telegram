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

	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args, " "))
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

	result, _ := externalCommand.ProcessText("!external_module arg1 arg2 arg3", userTest, "testchat", false)
	assert.Equal(t, "/home/ssalvatori/module/external_module ssalvatori 100 testchat arg1 arg2 arg3\n", result, "external")

	_, err := externalCommand.ProcessText("!mock-error arg1 arg2", userTest, "testchat", false)
	// assert.Equal(t, "Internal error", err.Error(), "external mock")
	assert.Nil(t, err)

}

func fakeExecCommandFail(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessFail", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_FAIL=1"}
	return cmd
}

func TestHelperProcessFail(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_FAIL") != "1" {
		return
	}
	os.Exit(1)
}

func TestRunCommandError(t *testing.T) {
	ExecCommand = fakeExecCommandFail
	defer func() { ExecCommand = exec.Command }()

	out := externalCommand.RunCommand("some_command", "arg1")
	assert.Equal(t, "", out, "Run Command Error")
}

func TestExternalCommandInject(t *testing.T) {
	userTest.Level = 100

	result, _ := externalCommand.ProcessText("!../../test arg1 arg2 arg3", userTest, "testchat", false)

	assert.Equal(t, "", result, "external command inject path traversal")
}

func TestExternalCommandSanitizesControlChars(t *testing.T) {
	userTest.Level = 100

	LookPathCommand = fakeLookPathCommand
	ExecCommand = fakeExecCommand

	defer func() {
		ExecCommand = exec.Command
		LookPathCommand = exec.LookPath
	}()

	// Null bytes and newlines in args/chat must be stripped before exec.
	// "arg1\x00arg2\narg3" → "arg1arg2arg3"; "testchat\x00injected" → "testchatinjected"
	result, err := externalCommand.ProcessText("!external_module arg1\x00arg2\narg3", userTest, "testchat\x00injected", false)
	assert.Nil(t, err)
	assert.NotContains(t, result, "\x00", "null byte must be stripped from args")
	assert.Contains(t, result, "arg1arg2arg3", "args must appear with control chars removed")
	assert.Contains(t, result, "testchatinjected", "chat name must appear with null byte removed")
}

func TestExternalCommandSanitizesUsername(t *testing.T) {
	userTest.Level = 100
	originalUsername := userTest.Username
	userTest.Username = "evil\x00user\nname"

	defer func() {
		userTest.Username = originalUsername
		ExecCommand = exec.Command
		LookPathCommand = exec.LookPath
	}()

	LookPathCommand = fakeLookPathCommand
	ExecCommand = fakeExecCommand

	// "evil\x00user\nname" → "evilusername"
	result, err := externalCommand.ProcessText("!external_module somearg", userTest, "testchat", false)
	assert.Nil(t, err)
	assert.NotContains(t, result, "\x00", "null byte must be stripped from username")
	assert.Contains(t, result, "evilusername", "username must appear with control chars removed")
}
