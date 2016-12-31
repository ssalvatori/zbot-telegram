package command

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var externalCommand = ExternalCommand{
	PathModules: "../modules/",
}

func TestExternalCommandOK(t *testing.T) {
	assert.Equal(t, "OK arg1 arg2 arg3\n",externalCommand.ProcessText("!test arg1 arg2 arg3", user), "external")
	externalCommand.Next = &FakeCommand{}
	assert.Equal(t, "",externalCommand.ProcessText("!date arg1 arg2", user), "external")
}

func TestExternalCommandInject(t *testing.T) {
	assert.Equal(t, "",externalCommand.ProcessText("!../../test arg1 arg2 arg3", user), "external")
}
