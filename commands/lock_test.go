package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var lockCommand = LockCommand{}

func TestTemplateCommandOK(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "100",
	}

	result, _ := lockCommand.ProcessText("!lock foo", userTest)
	assert.Equal(t, "[foo] locked", result, "Template Command")
}

func TestTemplateCommandErro(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}

	_, err := lockCommand.ProcessText("!lock foo", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
