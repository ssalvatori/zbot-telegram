package command

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var findCommand = FindCommand{}

func TestFindCommandOK(t *testing.T) {

	findCommand.Db = &db.MockZbotDatabase{
		Term:    "bar",
		Meaning: "bar",
	}
	assert.Equal(t, "bar", findCommand.ProcessText("!find foo", user), "Last Command")
	findCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}
	assert.Equal(t, "", findCommand.ProcessText("!find", user), "Last no next command")
	findCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", findCommand.ProcessText("?? ", user), "Last next command")
}
