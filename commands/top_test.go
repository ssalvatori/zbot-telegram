package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var topCommand = TopCommand{}

func TestTopCommandOK(t *testing.T) {
	topCommand.Db = &db.ZbotDatabaseMock{
		Level:     "7",
		FindTerms: []string{"foo", "bar"},
		RandDef:   []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
	}

	result, _ := topCommand.ProcessText("!top", userTest, "testchat", false)
	assert.Equal(t, "foo bar", result, "Top Command")

}

func TestTopCommandNotMatch(t *testing.T) {

	result, _ := topCommand.ProcessText("!top6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := topCommand.ProcessText("!top6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestTopCommandError(t *testing.T) {

	topCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}
	_, err := topCommand.ProcessText("!top", userTest, "testchat", false)
	assert.Equal(t, "Internal error, check logs", err.Error(), "Db error")

	_, err = topCommand.ProcessText("!top", userTest, "testchat", false)
	assert.Error(t, err, "Private mesage")

}
