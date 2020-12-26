package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var versionCommand = VersionCommand{
	Version:   "0.1",
	BuildTime: "2017-04-16 11:25:17.626575284 +0300 EEST",
	GitHash:   "6fd28bf",
}

func TestVersionCommandOK(t *testing.T) {
	result, _ := versionCommand.ProcessText("!version", userTest, "testchat", false)
	assert.Equal(t, "zbot golang version [0.1] commit [6fd28bf] build-time [2017-04-16 11:25:17.626575284 +0300 EEST]", result, "version command OK")
}

func TestVersionCommandNotMatch(t *testing.T) {

	result, _ := statsCommand.ProcessText("!version6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := statsCommand.ProcessText("!version", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")

	_, err = statsCommand.ProcessText("!version", userTest, "testchat", true)
	assert.Error(t, err, "Private message")
}
