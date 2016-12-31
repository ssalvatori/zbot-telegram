package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var statsCommand = StatsCommand{}

func TestStatsCommandOK(t *testing.T) {
	statsCommand.Db = &db.MockZbotDatabase{
		Level: "7",
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}
	assert.Equal(t, "Count: 7", statsCommand.ProcessText("!stats"), "Stats Command")
	assert.Equal(t, "", statsCommand.ProcessText("!stats6"), "Stats no next command")
	statsCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", statsCommand.ProcessText("!stats6"), "Stats next command")
}
