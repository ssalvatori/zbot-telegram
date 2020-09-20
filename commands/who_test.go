package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var whoCommand = WhoCommand{}

func TestWhoCommand(t *testing.T) {

	whoCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
		Author:  "ssalvatori",
		Date:    "2017-03-22",
		Level:   "100",
	}
	result, _ := whoCommand.ProcessText("!who foo", userTest, "testchat")
	assert.Equal(t, "[foo] by [ssalvatori] on [0001-01-01] hits [0]", result, "Who Command OK")
}

func TestWhoCommandNotMatch(t *testing.T) {

	result, _ := whoCommand.ProcessText("!who6", userTest, "testchat")
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := whoCommand.ProcessText("!who6", userTest, "testchat")
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestWhoCommandError(t *testing.T) {

	whoCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := whoCommand.ProcessText("!who foo", userTest, "testchat")
	// assert.Equal(t, "Internal error", err.Error(), "Db error")
	assert.Error(t, err, "DB Error")
}
