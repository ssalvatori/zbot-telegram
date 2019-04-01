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
	result, _ := versionCommand.ProcessText("!version", userTest)
	assert.Equal(t, "zbot golang version [0.1] commit [6fd28bf] build-time [2017-04-16 11:25:17.626575284 +0300 EEST]", result, "version command OK")
}

func TestVersionCommandNoNext(t *testing.T) {
	result, _ := versionCommand.ProcessText("!version6", userTest)
	assert.Equal(t, "", result, "version command no next")
}

func TestVersionCommandNoMatch(t *testing.T) {
	_, err := versionCommand.ProcessText("!version6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Version doesn't match")
}

/*
func TestVersionCommandNext(t *testing.T) {
	versionCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", versionCommand.ProcessText("!version6", userTest), "version command no next")
}
*/
