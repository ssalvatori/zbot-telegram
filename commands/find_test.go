package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var findCommand = FindCommand{}

func TestFindCommandOK(t *testing.T) {

	findCommand.Db = &db.MockZbotDatabase{
		Term: "bar",
		Meaning: "bar",
	}
	assert.Equal(t, "bar", findCommand.ProcessText("!find foo"), "Last Command")
	findCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}
	assert.Equal(t, "", findCommand.ProcessText("!find"), "Last no next command")
	findCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", findCommand.ProcessText("?? "), "Last next command")
}
