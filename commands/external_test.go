package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var externalCommand = ExternalCommand{
	PathModules: "../modules/",
}

func TestExternalCommandOK(t *testing.T) {
	assert.Equal(t, "OK ssalvatori 5 arg1 arg2 arg3\n", externalCommand.ProcessText("!test arg1 arg2 arg3", userTest), "external")
	externalCommand.Next = &FakeCommand{}
	assert.Equal(t, "OK ssalvatori 5 arg1 arg2\n", externalCommand.ProcessText("!test arg1 arg2", userTest), "external")
}

func TestExternalCommandInject(t *testing.T) {
	externalCommand.Next = nil
	assert.Equal(t, "", externalCommand.ProcessText("!../../test arg1 arg2 arg3", userTest), "external")
}
