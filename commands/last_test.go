package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var lastCommand = LastCommand{}

func TestLastCommandOK(t *testing.T) {
	lastCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}
	assert.Equal(t, "[foo] - [bar]", lastCommand.ProcessText("!last", userTest), "Last Command")
	assert.Equal(t, "", lastCommand.ProcessText("!last6", userTest), "Last no next command")
	lastCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", lastCommand.ProcessText("!last6", userTest), "Last next command")
}
