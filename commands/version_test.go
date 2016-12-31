package command

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var versionCommand = VersionCommand{
	Version: "0.1",
}

func TestVersionCommandOK(t *testing.T) {
	assert.Equal(t, "zbot golang version 0.1", versionCommand.ProcessText("!version"), "version command OK")
}

func TestVersionCommandNoNext(t *testing.T) {
	assert.Equal(t, "", versionCommand.ProcessText("!version6"), "version command no next")
}

func TestVersionCommandNext(t *testing.T) {
	versionCommand.Next = &FakeCommand{}
	assert.Equal(t, "", versionCommand.ProcessText("!version6"), "version command no next")
}