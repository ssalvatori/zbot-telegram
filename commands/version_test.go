package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var versionCommand = VersionCommand{
	Version:   "0.1",
	BuildTime: "2017-04-16 11:25:17.626575284 +0300 EEST",
}

func TestVersionCommandOK(t *testing.T) {
	assert.Equal(t, "zbot golang version [0.1] build-time [2017-04-16 11:25:17.626575284 +0300 EEST]", versionCommand.ProcessText("!version", userTest), "version command OK")
}

func TestVersionCommandNoNext(t *testing.T) {
	assert.Equal(t, "", versionCommand.ProcessText("!version6", userTest), "version command no next")
}

func TestVersionCommandNext(t *testing.T) {
	versionCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", versionCommand.ProcessText("!version6", userTest), "version command no next")
}
