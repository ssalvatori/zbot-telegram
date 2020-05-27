package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var lastCommand = LastCommand{}

func TestLastCommandOK(t *testing.T) {
	lastCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}
	result, _ := lastCommand.ProcessText("!last", userTest)
	assert.Equal(t, "[foo] - [bar]", result, "Last Command")
}

func TestLastCommandNotMatch(t *testing.T) {
	result, _ := lastCommand.ProcessText("!last6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := lastCommand.ProcessText("!last6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestLastCommandError(t *testing.T) {
	lastCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Error:   true,
	}
	_, err := lastCommand.ProcessText("!last", userTest)
	assert.Equal(t, "mock", err.Error(), "Db Error")
}
