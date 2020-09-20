package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var lastCommand = LastCommand{}

func TestLastCommandOK(t *testing.T) {
	lastCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
	}
	result, _ := lastCommand.ProcessText("!last", userTest, "testchat")
	assert.Equal(t, "[ foo ]", result, "Last Command")
}

func TestLastCommandNotMatch(t *testing.T) {
	result, _ := lastCommand.ProcessText("!last6", userTest, "testchat")
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := lastCommand.ProcessText("!last6", userTest, "testchat")
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestLastCommandError(t *testing.T) {
	lastCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
		Error:   true,
	}
	_, err := lastCommand.ProcessText("!last", userTest, "testchat")
	// assert.Equal(t, "Internal error", err.Error(), "Db Error")
	assert.Error(t, err, "Internal Error")
}

func TestPrintTerms(t *testing.T) {
	var items = []db.Definition{
		{Term: "term1", Meaning: "meaning 1"},
		{Term: "term2", Meaning: "meaning 2"},
	}

	result := PrintTerms(items)
	assert.Equal(t, "[ term1 term2 ]", result, "Last Command")
}
