package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var randCommand = RandCommand{}

func TestRandCommandOK(t *testing.T) {
	randCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}
	assert.Equal(t, "[foo] - [bar]", randCommand.ProcessText("!rand", user), "Rand Command")
	assert.Equal(t, "", randCommand.ProcessText("!rand6", user), "Rand no next command")
	randCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", randCommand.ProcessText("!ping6", user), "Rand next command")
}
