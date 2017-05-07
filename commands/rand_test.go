package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var randCommand = RandCommand{}

func TestRandCommandOK(t *testing.T) {

	randCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}

	assert.Equal(t, "", randCommand.ProcessText("!rand", userTest), "Rand error handler")

	randCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}
	assert.Equal(t, "[foo] - [bar]", randCommand.ProcessText("!rand", userTest), "Rand Command")
	assert.Equal(t, "", randCommand.ProcessText("!rand6", userTest), "Rand no next command")

	randCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", randCommand.ProcessText("!ping6", userTest), "Rand next command")
}
