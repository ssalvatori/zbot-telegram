package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var learnCommand = LearnCommand{}

func TestLearnCommandOK(t *testing.T) {
	var result string
	learnCommand.Db = &db.ZbotDatabaseMock{}

	result, _ = learnCommand.ProcessText("!learn foo bar", userTest, "test", false)
	assert.Equal(t, "[foo] - [bar]", result, "Lean Command")

}

func TestLearnCommandNotMatch(t *testing.T) {

	result, err := learnCommand.ProcessText("!learn6 foor ala", userTest, "test", false)
	assert.Equal(t, "", result, "Empty output")
	assert.Error(t, err, "Command doesn't match")

}

func TestLearnCommandError(t *testing.T) {

	learnCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := learnCommand.ProcessText("!learn foo lala", userTest, "test", false)
	assert.Error(t, err, "Internal error")

	_, err = learnCommand.ProcessText("!learn foo lala", userTest, "test", true)
	assert.Error(t, err, "Private message")
}
