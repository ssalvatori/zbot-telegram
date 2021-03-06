package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var lockCommand = LockCommand{}

func TestTemplateCommandOK(t *testing.T) {

	lockCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
		Level:   "100",
	}

	result, _ := lockCommand.ProcessText("!lock foo", userTest, "testchat", false)
	assert.Equal(t, "[foo] locked", result, "Template Command")
}

func TestTemplateCommandErro(t *testing.T) {

	lockCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}

	_, err := lockCommand.ProcessText("!lock foo", userTest, "testchat", false)
	assert.Error(t, err, "Internal error")

	_, err = lockCommand.ProcessText("!lock foo", userTest, "testchat", true)
	assert.Error(t, err, "Private message")
}
