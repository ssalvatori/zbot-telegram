package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
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

}

func TestStatsCommandNotMatch(t *testing.T) {

	result, _ := statsCommand.ProcessText("!stats6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := statsCommand.ProcessText("!stats6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestStatsCommandError(t *testing.T) {

	statsCommand.Db = &db.MockZbotDatabase{
		Level:    "7",
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := statsCommand.ProcessText("!stats", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
