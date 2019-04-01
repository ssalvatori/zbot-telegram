package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var statsCommand = StatsCommand{}

func TestStatsCommandOK(t *testing.T) {
	statsCommand.Db = &db.MockZbotDatabase{
		Level:    "7",
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}

	result, _ := statsCommand.ProcessText("!stats", userTest)
	assert.Equal(t, "Count: 7", result, "Stats Command")

	result, _ = statsCommand.ProcessText("!stats6", userTest)
	assert.Equal(t, "", result, "Stats no next command")

	_, err := statsCommand.ProcessText("!stats6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Stats next command")

	statsCommand.Db = &db.MockZbotDatabase{
		Level:    "7",
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err = statsCommand.ProcessText("!stats", userTest)
	assert.Equal(t, "mock", err.Error(), "Stats next command")
}
