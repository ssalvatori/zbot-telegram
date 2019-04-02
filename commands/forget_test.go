package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var forgetCommand = ForgetCommand{}

func TestForgetCommandOK(t *testing.T) {

	forgetCommand.Db = &db.MockZbotDatabase{}

	result, _ := forgetCommand.ProcessText("!forget foo", userTest)

	assert.Equal(t, "[foo] deleted", result, "Forget Command OK")
}

func TestForgetCommandNotMatch(t *testing.T) {

	result, _ := forgetCommand.ProcessText("!forget6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := forgetCommand.ProcessText("!forget6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestForgetCommandError(t *testing.T) {

	forgetCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}
	_, err := forgetCommand.ProcessText("!forget lal", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
