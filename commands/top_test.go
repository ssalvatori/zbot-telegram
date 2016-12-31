package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var topCommand = TopCommand{}

func TestTopCommandOK(t *testing.T) {
	topCommand.Db = &db.MockZbotDatabase{
		Level: "7",
		Find_terms: []string{"foo", "bar",},
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}
	assert.Equal(t, "foo bar", topCommand.ProcessText("!top", user), "Top Command")
	assert.Equal(t, "", topCommand.ProcessText("!top6", user), "Top no next command")
	topCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", topCommand.ProcessText("!top6", user), "Top next command")
}
