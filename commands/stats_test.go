package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var statsCommand = StatsCommand{}

func TestStatsCommandOK(t *testing.T) {
	statsCommand.Db = &db.ZbotDatabaseMock{
		Level:   "7",
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
	}

	result, _ := statsCommand.ProcessText("!stats", userTest, "testchat", false)
	assert.Equal(t, "Number of definitions: 7", result, "Stats Command")

}

func TestStatsCommandNotMatch(t *testing.T) {

	result, _ := statsCommand.ProcessText("!stats6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := statsCommand.ProcessText("!stats6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestStatsCommandError(t *testing.T) {

	statsCommand.Db = &db.ZbotDatabaseMock{
		Level:   "7",
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := statsCommand.ProcessText("!stats", userTest, "testchat", false)
	// assert.Equal(t, "Internal error", err.Error(), "Db error")
	assert.Error(t, err, "Internal error")
}
