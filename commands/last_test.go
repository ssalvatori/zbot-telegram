package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var lastCommand = LastCommand{}

func TestLastCommandOK(t *testing.T) {
	lastCommand.Db = &db.MockZbotDatabase{
		Term: "foo",
		Meaning: "bar",
	}
	assert.Equal(t, "[foo] - [bar]", lastCommand.ProcessText("!last"), "Last Command")
	assert.Equal(t, "", lastCommand.ProcessText("!last6"), "Last no next command")
	lastCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", lastCommand.ProcessText("!last6"), "Last next command")
}
