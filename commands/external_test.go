package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var externalCommand = ExternalCommand{
	PathModules: "../modules/",
}

func TestExternalCommandOK(t *testing.T) {
	//result, _ := externalCommand.ProcessText("!test arg1 arg2 arg3", userTest)
	//assert.Equal(t, "OK ssalvatori 5 arg1 arg2 arg3\n", result, "external")

	//assert.Equal(t, "OK ssalvatori 5 arg1 arg2\n", externalCommand.ProcessText("!test arg1 arg2", userTest), "external")
}

func TestExternalCommandInject(t *testing.T) {

	result, _ := externalCommand.ProcessText("!../../test arg1 arg2 arg3", userTest)

	assert.Equal(t, "", result, "external commmand inject")
}
