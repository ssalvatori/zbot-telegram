package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var forgetCommand = ForgetCommand{}

func TestForgetCommandOK(t *testing.T) {

	forgetCommand.Db = &db.ZbotDatabaseMock{}

	result, _ := forgetCommand.ProcessText("!forget foo", userTest, "testchat", false)

	assert.Equal(t, "[foo] deleted", result, "Forget Command OK")
}

func TestForgetCommandNotMatch(t *testing.T) {

	result, _ := forgetCommand.ProcessText("!forget6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := forgetCommand.ProcessText("!forget6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestForgetCommandError(t *testing.T) {

	forgetCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}
	_, err := forgetCommand.ProcessText("!forget lal", userTest, "testchat", false)
	assert.Error(t, err, "DB error")
}
